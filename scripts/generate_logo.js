const fs = require('fs');
const path = require('path');
const { Resvg } = require('B:/workgit/agentspec/node_modules/@resvg/resvg-js');

function buildKingfisherLogoSvg() {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" width="1024" height="1024">
  <defs>
    <clipPath id="squircle-clip">
      <rect x="24" y="24" width="976" height="976" rx="220" />
    </clipPath>

    <linearGradient id="cyan-glow" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#00f5ff"/>
      <stop offset="100%" stop-color="#0284c7"/>
    </linearGradient>

    <linearGradient id="cyan-glow-rev" x1="100%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#00f5ff"/>
      <stop offset="100%" stop-color="#0284c7"/>
    </linearGradient>

    <linearGradient id="amber-accent" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#f59e0b"/>
      <stop offset="100%" stop-color="#ea580c"/>
    </linearGradient>

    <linearGradient id="beak-blade" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#00f5ff"/>
      <stop offset="50%" stop-color="#0284c7"/>
      <stop offset="100%" stop-color="#0f172a"/>
    </linearGradient>

    <filter id="subtle-shadow" x="-10%" y="-10%" width="120%" height="120%">
      <feDropShadow dx="0" dy="16" stdDeviation="20" flood-color="#000000" flood-opacity="0.16" />
    </filter>
  </defs>

  <!-- Luxury White Squircle Container -->
  <rect x="24" y="24" width="976" height="976" rx="220" fill="#ffffff" stroke="#e2e8f0" stroke-width="6" />

  <g clip-path="url(#squircle-clip)">
    <g transform="translate(512, 512)" filter="url(#subtle-shadow)">

      <!-- Hexagonal Architectural Gateway -->
      <polygon points="
        0,-390
        338,-195
        338,195
        0,390
        -338,195
        -338,-195
      " fill="none" stroke="#0f172a" stroke-width="40" stroke-linejoin="round" />

      <polygon points="
        0,-355
        307,-177
        307,177
        0,355
        -307,177
        -307,-177
      " fill="none" stroke="#00f5ff" stroke-width="4" opacity="0.4" stroke-dasharray="16, 12" />

      <!-- Symmetrical Origami Mascot: Halcyon Kingfisher (Sub-15ms Dynamic Vector Craftsman) -->

      <!-- 1. Watertight Solid Base Silhouette -->
      <path d="
        M 0 -290
        L 75 -240
        L 120 -280
        L 115 -180
        L 190 -100
        L 215 10
        L 180 150
        L 110 240
        L 0 340
        L -110 240
        L -180 150
        L -215 10
        L -190 -100
        L -115 -180
        L -120 -280
        L -75 -240
        Z
      " fill="#0b0f19" />

      <!-- 2. Crest & Crown Armor (Polygonal Facets) -->
      <!-- Left Crest Feathers -->
      <polygon points="0,-290 -75,-240 -120,-280 -95,-190" fill="#1e293b" />
      <polygon points="-120,-280 -115,-180 -95,-190" fill="#334155" />
      <!-- Right Crest Feathers (Symmetrical) -->
      <polygon points="0,-290 75,-240 120,-280 95,-190" fill="#0f172a" />
      <polygon points="120,-280 115,-180 95,-190" fill="#1e293b" />

      <!-- Center Crown Ridge -->
      <polygon points="0,-290 -50,-210 0,-175 50,-210" fill="#334155" />
      <polygon points="0,-175 -50,-210 -75,-240 0,-290" fill="#1e293b" />
      <polygon points="0,-175 50,-210 75,-240 0,-290" fill="#475569" />

      <!-- 3. Upper Brow & Orbital Armor -->
      <!-- Left Brow -->
      <polygon points="0,-175 -50,-210 -110,-120 -40,-100 0,-120" fill="#1e293b" />
      <polygon points="-40,-100 -110,-120 -140,-30 -50,-40 0,-75" fill="#0f172a" />
      <!-- Right Brow (Symmetrical) -->
      <polygon points="0,-175 50,-210 110,-120 40,-100 0,-120" fill="#334155" />
      <polygon points="40,-100 110,-120 140,-30 50,-40 0,-75" fill="#1e293b" />

      <!-- 4. Electric Cyan Optics (Intelligent Predator Eyes) -->
      <!-- Left Eye Socket -->
      <polygon points="-35,-110 -85,-110 -105,-65 -70,-50 -30,-65" fill="#0b0f19" stroke="#00f5ff" stroke-width="2" />
      <polygon points="-42,-95 -75,-95 -90,-70 -68,-60 -38,-70" fill="url(#cyan-glow)" />
      <polygon points="-58,-90 -70,-70 -58,-65" fill="#ffffff" />
      <!-- Right Eye Socket (Symmetrical) -->
      <polygon points="35,-110 85,-110 105,-65 70,-50 30,-65" fill="#0b0f19" stroke="#00f5ff" stroke-width="2" />
      <polygon points="42,-95 75,-95 90,-70 68,-60 38,-70" fill="url(#cyan-glow-rev)" />
      <polygon points="58,-90 70,-70 58,-65" fill="#ffffff" />

      <!-- 5. Cheek Auricular Flashes (Kingfisher Cyan/Amber Plumage) -->
      <!-- Left Cheek Flashes -->
      <polygon points="-110,-120 -190,-100 -140,-30" fill="url(#cyan-glow)" />
      <polygon points="-140,-30 -190,-100 -215,10 -150,50" fill="#1e293b" />
      <polygon points="-140,-30 -150,50 -60,40 -50,-40" fill="url(#amber-accent)" />
      <!-- Right Cheek Flashes (Symmetrical) -->
      <polygon points="110,-120 190,-100 140,-30" fill="url(#cyan-glow-rev)" />
      <polygon points="140,-30 190,-100 215,10 150,50" fill="#0f172a" />
      <polygon points="140,-30 150,50 60,40 50,-40" fill="url(#amber-accent)" />

      <!-- 6. Aerodynamic Dagger Beak (Pinpoint Vector Precision) -->
      <!-- Upper Mandible Center Spine -->
      <polygon points="0,-75 -25,-10 0,160" fill="#00f5ff" />
      <polygon points="0,-75 25,-10 0,160" fill="#0284c7" />
      <!-- Beak Flanks -->
      <polygon points="-25,-10 -50,-40 -40,30 0,160" fill="#0f172a" />
      <polygon points="25,-10 50,-40 40,30 0,160" fill="#1e293b" />
      <!-- Lower Mandible Tip Blade -->
      <polygon points="0,160 -18,175 0,225" fill="#00f5ff" />
      <polygon points="0,160 18,175 0,225" fill="#0284c7" />

      <!-- 7. Throat Shield & Mantle Armor Plates -->
      <polygon points="-40,30 -60,40 -90,120 0,180" fill="#1e293b" />
      <polygon points="40,30 60,40 90,120 0,180" fill="#334155" />
      <polygon points="-90,120 -150,50 -180,150 -80,210 0,180" fill="#0f172a" />
      <polygon points="90,120 150,50 180,150 80,210 0,180" fill="#1e293b" />

      <!-- 8. Breastplate & Kinetic Vector Chevron -->
      <polygon points="0,180 -80,210 -110,240 0,285" fill="#1e293b" />
      <polygon points="0,180 80,210 110,240 0,285" fill="#334155" />
      <polygon points="0,285 -110,240 0,340" fill="url(#cyan-glow)" />
      <polygon points="0,285 110,240 0,340" fill="url(#cyan-glow-rev)" />

      <!-- Center Telemetry Pulse Ring -->
      <circle cx="0" cy="-75" r="5" fill="#00f5ff" />
      <circle cx="0" cy="180" r="4" fill="#00f5ff" />

    </g>
  </g>
</svg>`;
}

async function renderLogo() {
  const outputDir = path.resolve('B:/workgit/ogcraft/docs/images');
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  const svg = buildKingfisherLogoSvg();
  const svgPath = path.join(outputDir, 'logo.svg');
  const pngPath = path.join(outputDir, 'logo.png');

  fs.writeFileSync(svgPath, svg, 'utf-8');
  console.log('Saved SVG to:', svgPath);

  const resvg = new Resvg(svg, {
    fitTo: { mode: 'width', value: 1024 },
  });
  const pngData = resvg.render().asPng();
  fs.writeFileSync(pngPath, pngData);
  console.log('Saved 1024x1024 PNG to:', pngPath);
}

renderLogo().catch(console.error);
