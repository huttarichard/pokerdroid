import React from "react";
import C2 from "./C2";
import C3 from "./C3";
import C4 from "./C4";
import C5 from "./C5";
import C6 from "./C6";
import C7 from "./C7";
import C8 from "./C8";
import C9 from "./C9";
import CT from "./CT";
import CJ from "./CJ";
import CQ from "./CQ";
import CK from "./CK";
import CA from "./CA";
import D2 from "./D2";
import D3 from "./D3";
import D4 from "./D4";
import D5 from "./D5";
import D6 from "./D6";
import D7 from "./D7";
import D8 from "./D8";
import D9 from "./D9";
import DT from "./DT";
import DJ from "./DJ";
import DQ from "./DQ";
import DK from "./DK";
import DA from "./DA";
import H2 from "./H2";
import H3 from "./H3";
import H4 from "./H4";
import H5 from "./H5";
import H6 from "./H6";
import H7 from "./H7";
import H8 from "./H8";
import H9 from "./H9";
import HT from "./HT";
import HJ from "./HJ";
import HQ from "./HQ";
import HK from "./HK";
import HA from "./HA";
import S2 from "./S2";
import S3 from "./S3";
import S4 from "./S4";
import S5 from "./S5";
import S6 from "./S6";
import S7 from "./S7";
import S8 from "./S8";
import S9 from "./S9";
import ST from "./ST";
import SJ from "./SJ";
import SQ from "./SQ";
import SK from "./SK";
import SA from "./SA";
import B1 from "./B1";
import B2 from "./B2";

interface CardProps {
  card: string;
  height?: string;
  width?: string;
}

const cardComponents: Record<string, React.ComponentType<any>> = {
  "2c": C2,
  "3c": C3,
  "4c": C4,
  "5c": C5,
  "6c": C6,
  "7c": C7,
  "8c": C8,
  "9c": C9,
  tc: CT,
  jc: CJ,
  qc: CQ,
  kc: CK,
  ac: CA,
  "2d": D2,
  "3d": D3,
  "4d": D4,
  "5d": D5,
  "6d": D6,
  "7d": D7,
  "8d": D8,
  "9d": D9,
  td: DT,
  jd: DJ,
  qd: DQ,
  kd: DK,
  ad: DA,
  "2h": H2,
  "3h": H3,
  "4h": H4,
  "5h": H5,
  "6h": H6,
  "7h": H7,
  "8h": H8,
  "9h": H9,
  th: HT,
  jh: HJ,
  qh: HQ,
  kh: HK,
  ah: HA,
  "2s": S2,
  "3s": S3,
  "4s": S4,
  "5s": S5,
  "6s": S6,
  "7s": S7,
  "8s": S8,
  "9s": S9,
  ts: ST,
  js: SJ,
  qs: SQ,
  ks: SK,
  as: SA,
  back: B1,
  back2: B2,
};

export default function Card({
  card,
  height = "100px",
  width = "auto",
}: CardProps) {
  const CardComponent = cardComponents[card.toLowerCase()] || B1;

  return (
    <CardComponent
      style={{
        height,
        width,
        borderRadius: "4px",
        boxShadow: "0 2px 4px rgba(0,0,0,0.2)",
      }}
    />
  );
}
