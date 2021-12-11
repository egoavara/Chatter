/**
 * plugins/webfontloader.js
 *
 * webfontloader documentation: https://github.com/typekit/webfontloader
 */

import webfontloader from "webfontloader";

export async function loadFonts() {
  // const webFontLoader = webfontloader;

  webfontloader.load({
    google: {
      families: ["Roboto:100,300,400,500,700,900&display=swap"],
    },
  });
}
