import React from 'react';

export default function DA(props: any) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" xmlnsXlink="http://www.w3.org/1999/xlink" className="card" face="AD" height="3.5in" preserveAspectRatio="none" viewBox="-120 -168 240 336" width="2.5in" {...props}>
      <defs>
        <symbol id="VDA" viewBox="-500 -500 1000 1000" preserveAspectRatio="xMinYMid">
          <path d="M-270 460L-110 460M-200 450L0 -460L200 450M110 460L270 460M-120 130L120 130" stroke="red" strokeWidth={80} strokeLinecap="square" strokeMiterlimit="1.5" fill="none" />
        </symbol>
        <symbol id="SDA" viewBox="-600 -600 1200 1200" preserveAspectRatio="xMinYMid">
          <path d="M-400 0C-350 0 0 -450 0 -500C0 -450 350 0 400 0C350 0 0 450 0 500C0 450 -350 0 -400 0Z" fill="red" />
        </symbol>
      </defs>
      <rect width={239} height={335} x="-119.5" y="-167.5" rx={12} ry={12} fill="white" stroke="black" />
      <use xlinkHref="#SDA" height={40} width={40} x={-20} y={-20} />
      <use xlinkHref="#VDA" height={70} width={70} x={-122} y={-156} />
      <use xlinkHref="#SDA" height="58.558" width="58.558" x="-116.279" y={-81} />
      <g transform="rotate(180)">
        <use xlinkHref="#VDA" height={70} width={70} x={-122} y={-156} />
        <use xlinkHref="#SDA" height="58.558" width="58.558" x="-116.279" y={-81} />
      </g>
    </svg>
  );
}
