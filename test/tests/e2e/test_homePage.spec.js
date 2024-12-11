// @ts-check
import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveTitle(/HNG Boilerplate/);
});

test('get started link', async ({ page }) => {
  await page.goto('/');
  await page.getByRole('link', { name: 'Get Started' }).first().click();
  await expect(page.getByRole('heading', { name: 'Sign Up' })).toBeVisible();
});
