import { expect, test } from '@playwright/test'

test('real browser → Vite proxy → Go: arithmetic, validation, recovery', async ({ page }) => {
  const failures: string[] = []
  page.on('pageerror', error => failures.push(error.message))
  await page.goto('/')
  const expression = page.getByRole('textbox', { name: 'EXPRESSION' })
  for (const [input, expected] of [
    ['5+3', '8'], ['5-3', '2'], ['5*3', '15'], ['10/4', '2.5'],
    ['5+3*2', '11'], ['10-6/2', '7'], ['-3*2', '-6'], ['2*-3', '-6'], ['0.1+0.2', '0.3'],
  ]) {
    await expression.fill(input!)
    const response = page.waitForResponse('/api/calculate')
    await expression.press('Enter')
    expect((await response).status()).toBe(200)
    await expect(page.getByLabel('Result', { exact: true })).toHaveText(expected!)
  }
  for (const [input, expected] of [['5/0', 'Cannot divide by zero'], ['5+*2', 'Enter a valid expression']]) {
    await expression.fill(input!)
    const response = page.waitForResponse('/api/calculate')
    await expression.press('Enter')
    expect((await response).status()).toBe(400)
    await expect(page.getByRole('alert')).toContainText(expected!)
  }
  await page.getByRole('button', { name: 'Clear', exact: true }).click()
  for (const name of ['5', 'Add', '3', 'Multiply', '2', 'Calculate']) await page.getByRole('button', { name, exact: true }).click()
  await expect(page.getByLabel('Result', { exact: true })).toHaveText('11')
  await expect(page.getByRole('alert')).toBeEmpty()
  await page.getByRole('button', { name: 'Clear', exact: true }).click()
  await page.locator('h1').click()
  await page.keyboard.type('10/4')
  await page.keyboard.press('Enter')
  await expect(page.getByLabel('Result', { exact: true })).toHaveText('2.5')
  await page.keyboard.press('Escape')
  await expect(expression).toHaveValue('')
  expect(failures).toEqual([])
})

test('responsive layout and focus at required widths', async ({ page }) => {
  await page.goto('/')
  for (const width of [1440, 768, 375]) {
    await page.setViewportSize({ width, height: 900 })
    await page.getByRole('textbox').focus()
    const focus = await page.getByRole('textbox').evaluate(element => {
      const label = document.querySelector('label[for="expression"]')
      if (!label) throw new Error('Missing expression label')
      const style = getComputedStyle(element)
      return {
        top: element.getBoundingClientRect().top - parseFloat(style.outlineWidth) - parseFloat(style.outlineOffset),
        labelBottom: label.getBoundingClientRect().bottom,
        outline: style.outlineStyle,
      }
    })
    expect(focus.outline).not.toBe('none')
    expect(focus.top).toBeGreaterThanOrEqual(focus.labelBottom)
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
    await expect(page.getByRole('button', { name: 'Calculate' })).toBeInViewport()
    await page.getByRole('textbox').fill('1e308*1')
    await page.getByRole('textbox').press('Enter')
    await expect(page.getByLabel('Result', { exact: true })).toHaveText('1e+308')
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
    await page.screenshot({ path: `test-results/calculator-${width}.png`, fullPage: true })
  }
  await page.goto('/')
  await page.keyboard.press('Tab')
  await expect(page.getByRole('textbox')).toBeFocused()
  for (const name of ['Open parenthesis', 'Close parenthesis', 'Power', 'Square root', 'Percentage', 'Clear']) {
    await page.keyboard.press('Tab')
    await expect(page.getByRole('button', { name, exact: true })).toBeFocused()
  }
  expect(await page.getByRole('button', { name: 'Clear', exact: true }).evaluate(element => getComputedStyle(element).outlineStyle)).not.toBe('none')
})

test('advanced operations through real Go and new controls', async ({ page }) => {
  await page.goto('/')
  const expression = page.getByRole('textbox')
  for (const [input, expected] of [
    ['2^3', '8'], ['2^3^2', '512'], ['-2^2', '-4'], ['(-2)^2', '4'], ['2^-3', '0.125'],
    ['sqrt(81)', '9'], ['sqrt(2)', '1.41421356237'], ['2^3+sqrt(16)', '12'],
    ['200+10%', '220'], ['200-10%', '180'], ['200*10%', '20'], ['200/10%', '2000'],
    ['200+10%+5%', '231'], ['200+10%*2', '200.2'], ['sqrt(25)*20%', '1'], ['(2+3)*4','20'],
  ]) {
    await expression.fill(input!)
    await expression.press('Enter')
    await expect(page.getByLabel('Result', { exact:true })).toHaveText(expected!)
  }
  for (const input of ['sqrt(-1)', '(-2)^0.5', '2^', 'sqrt()', '200/0%']) {
    await expression.fill(input)
    const response = page.waitForResponse('/api/calculate')
    await expression.press('Enter')
    expect((await response).status()).toBe(400)
    await expect(page.getByRole('alert')).not.toBeEmpty()
  }
  await page.getByRole('button',{name:'Clear',exact:true}).click()
  for (const name of ['Square root','9','Close parenthesis','Add','Open parenthesis','2','Power','3','Close parenthesis','Multiply','2','0','Percentage','Calculate']) {
    await page.getByRole('button',{name,exact:true}).click()
  }
  await expect(expression).toHaveValue('sqrt(9)+(2^3)*20%')
  await expect(page.getByLabel('Result',{exact:true})).toHaveText('4.6')
  await page.getByRole('button',{name:'Clear',exact:true}).click()
  await page.locator('h1').click()
  await page.keyboard.type('sqrt(9)+2^3+10%')
  await page.keyboard.press('Enter')
  await expect(page.getByLabel('Result',{exact:true})).toHaveText('12.1')
})

test('transport failure is visible and recoverable (fault injection)', async ({ page }) => {
  await page.goto('/')
  await page.route('**/api/calculate', route => route.abort('connectionrefused'))
  await page.getByRole('textbox').fill('2+3')
  await page.getByRole('button', { name: 'Calculate' }).click()
  await expect(page.getByRole('alert')).toContainText('Cannot reach the calculator')
  await page.unroute('**/api/calculate')
  await page.getByRole('button', { name: 'Calculate' }).click()
  await expect(page.getByLabel('Result', { exact: true })).toHaveText('5')
})

test('a negative returned result remains the complete base of a continued power', async ({ page }) => {
  await page.goto('/')
  const expression = page.getByRole('textbox')
  for (const [exponent, expected] of [['2', '4'], ['0.5', 'domain_error']]) {
    await expression.fill('1-3')
    await expression.press('Enter')
    await expect(page.getByLabel('Result', {exact:true})).toHaveText('-2')
    await page.getByRole('button',{name:'Power',exact:true}).click()
    await page.keyboard.type(exponent!)
    await expect(expression).toHaveValue(`(-2)^${exponent}`)
    const response = page.waitForResponse('/api/calculate')
    await page.keyboard.press('Enter')
    const wire = await response
    if (expected === 'domain_error') {
      expect(wire.status()).toBe(400)
      expect((await wire.json() as {error:{code:string}}).error.code).toBe('domain_error')
      await expect(page.getByRole('alert')).toContainText('real-number domain')
    } else {
      expect(wire.status()).toBe(200)
      await expect(page.getByLabel('Result',{exact:true})).toHaveText(expected!)
    }
  }
})

test('stopped backend produces visible HTTP failure', async ({ page }) => {
  test.skip(process.env.CALCULATOR_BACKEND_DOWN !== '1', 'Run deliberately with Go stopped and Vite running')
  await page.goto('/')
  await page.getByRole('textbox').fill('2+3')
  const response = page.waitForResponse('/api/calculate')
  await page.getByRole('button', { name: 'Calculate' }).click()
  expect((await response).status()).toBe(502)
  await expect(page.getByRole('alert')).toContainText('The server returned an invalid response')
  await expect(page.getByRole('button', { name: 'Calculate' })).toBeEnabled()
})
