const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const puppeteer = require('B:/workgit/bundlelens/node_modules/puppeteer');

async function run() {
  const imagesDir = path.resolve('B:/workgit/ogcraft/docs/images');
  const framesDir = path.resolve('B:/workgit/ogcraft/temp_frames');
  if (!fs.existsSync(imagesDir)) fs.mkdirSync(imagesDir, { recursive: true });
  if (!fs.existsSync(framesDir)) fs.mkdirSync(framesDir, { recursive: true });

  console.log('Launching browser...');
  const browser = await puppeteer.launch({
    headless: 'new',
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--window-size=1440,900'],
  });

  const page = await browser.newPage();
  await page.setViewport({ width: 1440, height: 900, deviceScaleFactor: 2 });

  console.log('Navigating to http://localhost:8092...');
  await page.goto('http://localhost:8092', { waitUntil: 'networkidle0' });

  // Wait 1s for initial render
  await new Promise((r) => setTimeout(r, 1000));

  // 1. Desktop Studio Screenshot
  console.log('Capturing Desktop Studio Screenshot...');
  await page.screenshot({
    path: path.join(imagesDir, 'ogcraft_dashboard.png'),
    clip: { x: 0, y: 0, width: 1440, height: 900 },
  });
  await page.screenshot({
    path: path.join(framesDir, 'frame_01.png'),
    clip: { x: 0, y: 0, width: 1440, height: 900 },
  });

  // 2. Select "glow_dev" template for frame 2
  console.log('Switching to glow_dev template...');
  const glowItem = await page.evaluateHandle(() => {
    const items = Array.from(document.querySelectorAll('span'));
    return items.find((el) => el.textContent && el.textContent.includes('Glow Developer')) || null;
  });
  if (glowItem && glowItem.asElement()) {
    await glowItem.asElement().click();
    await new Promise((r) => setTimeout(r, 1200));
    await page.screenshot({
      path: path.join(framesDir, 'frame_02.png'),
      clip: { x: 0, y: 0, width: 1440, height: 900 },
    });
  }

  // 3. Select "product_launch" template for frame 3
  console.log('Switching to product_launch template...');
  const productItem = await page.evaluateHandle(() => {
    const items = Array.from(document.querySelectorAll('span'));
    return items.find((el) => el.textContent && el.textContent.includes('Product Launch')) || null;
  });
  if (productItem && productItem.asElement()) {
    await productItem.asElement().click();
    await new Promise((r) => setTimeout(r, 1200));
    await page.screenshot({
      path: path.join(framesDir, 'frame_03.png'),
      clip: { x: 0, y: 0, width: 1440, height: 900 },
    });
  }

  // 4. Switch to Social Feeds Simulator
  console.log('Switching to Social Feeds Simulator...');
  const feedBtn = await page.evaluateHandle(() => {
    const buttons = Array.from(document.querySelectorAll('button'));
    return buttons.find((b) => b.textContent && b.textContent.includes('Social Feeds')) || null;
  });
  if (feedBtn && feedBtn.asElement()) {
    await feedBtn.asElement().click();
    await new Promise((r) => setTimeout(r, 1200));

    // Save Social Simulator Screenshot
    await page.screenshot({
      path: path.join(imagesDir, 'ogcraft_social_simulator.png'),
      clip: { x: 0, y: 0, width: 1440, height: 900 },
    });
    await page.screenshot({
      path: path.join(framesDir, 'frame_04.png'),
      clip: { x: 0, y: 0, width: 1440, height: 900 },
    });
  }

  // 5. Mobile Responsive View
  console.log('Capturing Mobile Responsive View (390x844)...');
  const mobilePage = await browser.newPage();
  await mobilePage.setViewport({ width: 390, height: 844, deviceScaleFactor: 2, isMobile: true, hasTouch: true });
  await mobilePage.goto('http://localhost:8092', { waitUntil: 'networkidle0' });
  await new Promise((r) => setTimeout(r, 1000));

  await mobilePage.screenshot({
    path: path.join(imagesDir, 'ogcraft_mobile.png'),
    clip: { x: 0, y: 0, width: 390, height: 844 },
  });
  await mobilePage.close();

  await browser.close();
  console.log('Browser closed.');

  // Generate animated GIF from frames using ffmpeg
  console.log('Generating animated GIF with ffmpeg...');
  const gifPath = path.join(imagesDir, 'ogcraft_demo.gif');
  const cmd = `ffmpeg -y -framerate 1 -i "${path.join(framesDir, 'frame_%02d.png')}" -vf "scale=1200:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=128:reserve_transparent=0[p];[s1][p]paletteuse=dither=bayer:bayer_scale=3" "${gifPath}"`;
  execSync(cmd, { stdio: 'inherit' });

  // Clean up frames directory
  fs.rmSync(framesDir, { recursive: true, force: true });
  console.log('✓ Successfully created real screenshots and animated demo GIF at docs/images/');
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
