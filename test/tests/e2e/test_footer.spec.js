import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/'); 
  await expect(page).toHaveTitle(/HNG Boilerplate/);
});

test('About Us link', async ({ page }) => {
  await page.goto('/'); 
  const aboutUsLink = page.locator('text=About Us');
  const footer = page.locator('footer'); 
  await aboutUsLink.click();
  await expect(footer.locator('text=About Us')).toBeVisible();
});

test('FAQ link', async ({ page }) => {
  await page.goto('/');
  const footer = page.locator('footer');
  await footer.scrollIntoViewIfNeeded();
  const faqLink = footer.locator('text=FAQ');
  if (await faqLink.isVisible()) {
    await faqLink.click();
    await page.waitForNavigation();
    await expect(page).toHaveURL(/faqs/);
    await expect(page.locator('p', { hasText: 'FAQS' })).toBeVisible();
    await expect(footer.locator('text=FAQ')).toBeVisible();
  } else {
    throw new Error('FAQ link not found in the footer');
  }
});

test('Terms and Conditions link', async ({ page }) => {
  await page.goto('/'); // Use relative path
  const footer = page.locator('footer');
  await footer.scrollIntoViewIfNeeded();
  const termsLink = footer.locator('text=Terms and Condition');
  if (await termsLink.isVisible()) {
    await termsLink.click();
    await page.waitForNavigation();
    await expect(page).toHaveURL(/terms-and-conditions/);
    await expect(
      page.locator('p', { hasText: 'Terms and Conditions'})
    ).toBeVisible();
    await expect(footer.locator('text=Terms and Conditions')).toBeVisible();
  } else {
    throw new Error('Terms and Conditions link not found in the footer');
  }
});

test('Career link', async ({ page }) => {
  await page.goto('/'); 
  const footer = page.locator('footer');
  await footer.scrollIntoViewIfNeeded();
  const careerLink = footer.locator('text=Career');
  if (await careerLink.isVisible()) {
    await careerLink.click();
    await page.waitForNavigation();
    await expect(page).toHaveURL(/career/);
    await expect(
      page.locator('p').filter({ hasText: 'Career' }).first()
    ).toBeVisible();
    await expect(footer.locator('text=Career')).toBeVisible();
  } else {
    throw new Error('Career link not found in the footer');
  }
});

test('Waiting List link', async ({ page }) => {
  await page.goto('/');
  const footer = page.locator('footer');
  await footer.scrollIntoViewIfNeeded();
  const waitingListLink = footer.locator('text=Waiting List');
  if (await waitingListLink.isVisible()) {
    await waitingListLink.click();
    await page.waitForNavigation();
    await expect(page).toHaveURL(/waitlist/);
    await expect(
      page.locator('div').filter({ hasText: 'waitlist' }).first()
    ).toBeVisible();
    await expect(footer.locator('text=Waiting List')).toBeVisible();
  } else {
    throw new Error('Waitlist link not found in the footer');
  }
});

