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
  await page.keyboard.press('Tab')
  await expect(page.getByRole('button', { name: 'Clear', exact: true })).toBeFocused()
  expect(await page.getByRole('button', { name: 'Clear', exact: true }).evaluate(element => getComputedStyle(element).outlineStyle)).not.toBe('none')
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
