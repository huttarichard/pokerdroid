import React from 'react';

export default function S3(props: any) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" xmlnsXlink="http://www.w3.org/1999/xlink" className="card" face="3S" height="3.5in" preserveAspectRatio="none" viewBox="-120 -168 240 336" width="2.5in" {...props}>
      <defs>
        <symbol id="SS3" viewBox="-600 -600 1200 1200" preserveAspectRatio="xMinYMid">
          <path d="M0 -500C100 -250 355 -100 355 185A150 150 0 0 1 55 185A10 10 0 0 0 35 185C35 385 85 400 130 500L-130 500C-85 400 -35 385 -35 185A10 10 0 0 0 -55 185A150 150 0 0 1 -355 185C-355 -100 -100 -250 0 -500Z" fill="black" />
        </symbol>
        <symbol id="VS3" viewBox="-500 -500 1000 1000" preserveAspectRatio="xMinYMid">
          <path d="M-250 -320L-250 -460L200 -460L-110 -80C-100 -90 -50 -120 0 -120C200 -120 250 0 250 150C250 350 170 460 -30 460C-230 460 -260 300 -260 300" stroke="black" strokeWidth={80} strokeLinecap="square" strokeMiterlimit="1.5" fill="none" />
        </symbol>
      </defs>
      <rect width={239} height={335} x="-119.5" y="-167.5" rx={12} ry={12} fill="white" stroke="black" />
      <use xlinkHref="#VS3" height={70} width={70} x={-122} y={-156} />
      <use xlinkHref="#SS3" height="58.558" width="58.558" x="-116.279" y={-81} />
      <use xlinkHref="#SS3" height={40} width={40} x={-20} y="-107.668" />
      <use xlinkHref="#SS3" height={40} width={40} x={-20} y={-20} />
      <g transform="rotate(180)">
        <use xlinkHref="#VS3" height={70} width={70} x={-122} y={-156} />
        <use xlinkHref="#SS3" height="58.558" width="58.558" x="-116.279" y={-81} />
        <use xlinkHref="#SS3" height={40} width={40} x={-20} y="-107.668" />
      </g>
    </svg>
  );
}
