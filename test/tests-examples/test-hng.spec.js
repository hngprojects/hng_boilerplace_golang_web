import { test, expect } from '@playwright/test';

test('test', async ({ page }) => {
  await page.goto('https://kimiko-golang.teams.hng.tech/');
  await page
    .getByRole('navigation')
    .getByRole('link', { name: 'Home' })
    .click();
  await page
    .locator('div')
    .filter({ hasText: /^hero\.headlinehero\.descriptionhero\.cta$/ })
    .first()
    .click();
  await page.locator('.flex-1 > .bg-background').click();
  await page.locator('.swiper-slide > .flex').first().click();
  await page
    .locator('#swiper-wrapper-b2b95c96fbbb4fbf > div:nth-child(2) > .flex')
    .click();
  await page.locator('div:nth-child(3) > .flex').first().click();
  await page
    .locator('div')
    .filter({ hasText: 'Find The Perfect FitChoose' })
    .nth(3)
    .click();
  await page.getByRole('heading', { name: 'Find The Perfect Fit' }).click();
  await page.getByText('Choose the boilerplate plan').click();
  await page
    .locator('div')
    .filter({
      hasText: 'Boiler plateLogo subject details and addressSign Up For',
    })
    .nth(3)
    .click();
  await page.getByTestId('pricing').click();
  await page.getByTestId('pricing-tag').click();
  await page.getByText('Simple and').click();
  await page.getByText('Affordable').click();
  await page.getByText('Pricing Plan').click();
  await page.getByTestId('pricing-description').click();
  await page.getByTestId('monthly-toggle').click();
  await page.getByTestId('annual-toggle').click();
  await page
    .getByRole('button', { name: 'What is the purpose of this' })
    .click();
  await page
    .getByRole('button', { name: 'What is the purpose of this' })
    .click();
  await page
    .getByRole('button', { name: 'How do I reset my password?' })
    .click();
  await page
    .getByRole('button', { name: 'How do I reset my password?' })
    .click();
  await page
    .getByRole('button', { name: 'Can I use this application on' })
    .click();
  await page
    .getByRole('button', { name: 'Can I use this application on' })
    .click();
  await page
    .getByRole('heading', { name: 'Frequently Asked Questions' })
    .click();
  await page.getByTestId('contact-button').click();
});
