import { GlobalStyles } from "@mui/material";

import inter18ptBlack from "./fonts/Inter_18pt-Black.ttf";
import inter24ptBlack from "./fonts/Inter_24pt-Black.ttf";
import inter28ptBlack from "./fonts/Inter_28pt-Black.ttf";

import inter18ptBlackItalic from "./fonts/Inter_18pt-BlackItalic.ttf";
import inter24ptBlackItalic from "./fonts/Inter_24pt-BlackItalic.ttf";
import inter28ptBlackItalic from "./fonts/Inter_28pt-BlackItalic.ttf";

import inter18ptBold from "./fonts/Inter_18pt-Bold.ttf";
import inter24ptBold from "./fonts/Inter_24pt-Bold.ttf";
import inter28ptBold from "./fonts/Inter_28pt-Bold.ttf";

import inter18ptBoldItalic from "./fonts/Inter_18pt-BoldItalic.ttf";
import inter24ptBoldItalic from "./fonts/Inter_24pt-BoldItalic.ttf";
import inter28ptBoldItalic from "./fonts/Inter_28pt-BoldItalic.ttf";

import inter18ptExtraBold from "./fonts/Inter_18pt-ExtraBold.ttf";
import inter24ptExtraBold from "./fonts/Inter_24pt-ExtraBold.ttf";
import inter28ptExtraBold from "./fonts/Inter_28pt-ExtraBold.ttf";

import inter18ptExtraBoldItalic from "./fonts/Inter_18pt-ExtraBoldItalic.ttf";
import inter24ptExtraBoldItalic from "./fonts/Inter_24pt-ExtraBoldItalic.ttf";
import inter28ptExtraBoldItalic from "./fonts/Inter_28pt-ExtraBoldItalic.ttf";

import inter18ptExtraLight from "./fonts/Inter_18pt-ExtraLight.ttf";
import inter24ptExtraLight from "./fonts/Inter_24pt-ExtraLight.ttf";
import inter28ptExtraLight from "./fonts/Inter_28pt-ExtraLight.ttf";

import inter18ptExtraLightItalic from "./fonts/Inter_18pt-ExtraLightItalic.ttf";
import inter24ptExtraLightItalic from "./fonts/Inter_24pt-ExtraLightItalic.ttf";
import inter28ptExtraLightItalic from "./fonts/Inter_28pt-ExtraLightItalic.ttf";

import inter18ptItalic from "./fonts/Inter_18pt-Italic.ttf";
import inter24ptItalic from "./fonts/Inter_24pt-Italic.ttf";
import inter28ptItalic from "./fonts/Inter_28pt-Italic.ttf";

import inter18ptLight from "./fonts/Inter_18pt-Light.ttf";
import inter24ptLight from "./fonts/Inter_24pt-Light.ttf";
import inter28ptLight from "./fonts/Inter_28pt-Light.ttf";

import inter18ptLightItalic from "./fonts/Inter_18pt-LightItalic.ttf";
import inter24ptLightItalic from "./fonts/Inter_24pt-LightItalic.ttf";
import inter28ptLightItalic from "./fonts/Inter_28pt-LightItalic.ttf";

import inter18ptMedium from "./fonts/Inter_18pt-Medium.ttf";
import inter24ptMedium from "./fonts/Inter_24pt-Medium.ttf";
import inter28ptMedium from "./fonts/Inter_28pt-Medium.ttf";

import inter18ptMediumItalic from "./fonts/Inter_18pt-MediumItalic.ttf";
import inter24ptMediumItalic from "./fonts/Inter_24pt-MediumItalic.ttf";
import inter28ptMediumItalic from "./fonts/Inter_28pt-MediumItalic.ttf";

import inter18ptRegular from "./fonts/Inter_18pt-Regular.ttf";
import inter24ptRegular from "./fonts/Inter_24pt-Regular.ttf";
import inter28ptRegular from "./fonts/Inter_28pt-Regular.ttf";

import inter18ptSemiBold from "./fonts/Inter_18pt-SemiBold.ttf";
import inter24ptSemiBold from "./fonts/Inter_24pt-SemiBold.ttf";
import inter28ptSemiBold from "./fonts/Inter_28pt-SemiBold.ttf";

import inter18ptSemiBoldItalic from "./fonts/Inter_18pt-SemiBoldItalic.ttf";
import inter24ptSemiBoldItalic from "./fonts/Inter_24pt-SemiBoldItalic.ttf";
import inter28ptSemiBoldItalic from "./fonts/Inter_28pt-SemiBoldItalic.ttf";

import inter18ptThin from "./fonts/Inter_18pt-Thin.ttf";
import inter24ptThin from "./fonts/Inter_24pt-Thin.ttf";
import inter28ptThin from "./fonts/Inter_28pt-Thin.ttf";

import inter18ptThinItalic from "./fonts/Inter_18pt-ThinItalic.ttf";
import inter24ptThinItalic from "./fonts/Inter_24pt-ThinItalic.ttf";
import inter28ptThinItalic from "./fonts/Inter_28pt-ThinItalic.ttf";

// Font face declarations
const fontFaces = `
@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 100;
    font-display: swap;
    src: url('${inter18ptThin}') format('truetype'),
         url('${inter24ptThin}') format('truetype'),
         url('${inter28ptThin}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 100;
    font-display: swap;
    src: url('${inter18ptThinItalic}') format('truetype'),
         url('${inter24ptThinItalic}') format('truetype'),
         url('${inter28ptThinItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 200;
    font-display: swap;
    src: url('${inter18ptExtraLight}') format('truetype'),
         url('${inter24ptExtraLight}') format('truetype'),
         url('${inter28ptExtraLight}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 200;
    font-display: swap;
    src: url('${inter18ptExtraLightItalic}') format('truetype'),
         url('${inter24ptExtraLightItalic}') format('truetype'),
         url('${inter28ptExtraLightItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 300;
    font-display: swap;
    src: url('${inter18ptLight}') format('truetype'),
         url('${inter24ptLight}') format('truetype'),
         url('${inter28ptLight}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 300;
    font-display: swap;
    src: url('${inter18ptLightItalic}') format('truetype'),
         url('${inter24ptLightItalic}') format('truetype'),
         url('${inter28ptLightItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 400;
    font-display: swap;
    src: url('${inter18ptRegular}') format('truetype'),
         url('${inter24ptRegular}') format('truetype'),
         url('${inter28ptRegular}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 400;
    font-display: swap;
    src: url('${inter18ptItalic}') format('truetype'),
         url('${inter24ptItalic}') format('truetype'),
         url('${inter28ptItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 500;
    font-display: swap;
    src: url('${inter18ptMedium}') format('truetype'),
         url('${inter24ptMedium}') format('truetype'),
         url('${inter28ptMedium}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 500;
    font-display: swap;
    src: url('${inter18ptMediumItalic}') format('truetype'),
         url('${inter24ptMediumItalic}') format('truetype'),
         url('${inter28ptMediumItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 600;
    font-display: swap;
    src: url('${inter18ptSemiBold}') format('truetype'),
         url('${inter24ptSemiBold}') format('truetype'),
         url('${inter28ptSemiBold}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 600;
    font-display: swap;
    src: url('${inter18ptSemiBoldItalic}') format('truetype'),
         url('${inter24ptSemiBoldItalic}') format('truetype'),
         url('${inter28ptSemiBoldItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 700;
    font-display: swap;
    src: url('${inter18ptBold}') format('truetype'),
         url('${inter24ptBold}') format('truetype'),
         url('${inter28ptBold}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 700;
    font-display: swap;
    src: url('${inter18ptBoldItalic}') format('truetype'),
         url('${inter24ptBoldItalic}') format('truetype'),
         url('${inter28ptBoldItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 800;
    font-display: swap;
    src: url('${inter18ptExtraBold}') format('truetype'),
         url('${inter24ptExtraBold}') format('truetype'),
         url('${inter28ptExtraBold}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 800;
    font-display: swap;
    src: url('${inter18ptExtraBoldItalic}') format('truetype'),
         url('${inter24ptExtraBoldItalic}') format('truetype'),
         url('${inter28ptExtraBoldItalic}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: normal;
    font-weight: 900;
    font-display: swap;
    src: url('${inter18ptBlack}') format('truetype'),
         url('${inter24ptBlack}') format('truetype'),
         url('${inter28ptBlack}') format('truetype');
}

@font-face {
    font-family: 'Inter';
    font-style: italic;
    font-weight: 900;
    font-display: swap;
    src: url('${inter18ptBlackItalic}') format('truetype'),
         url('${inter24ptBlackItalic}') format('truetype'),
         url('${inter28ptBlackItalic}') format('truetype');
}
`;

export default function Inter() {
  return <GlobalStyles styles={fontFaces} />;
}
