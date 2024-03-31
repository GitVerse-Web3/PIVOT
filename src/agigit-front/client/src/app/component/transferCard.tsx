import { Card, CardBody, CardHeader, Code } from "@nextui-org/react";
import Image from "next/image";
import dottedArrow from "/public/imgs/dottedArrow.png";

export default function TransferCard(props: any) {
  const { result, extraData } = props;
  const formatAddress = (s: string) => {
    if (s) {
      const head = s.slice(0, 4);
      const tail = s.slice(s.length - 4, s.length);
      return head + "..." + tail;
    }
    return "...";
  };
  return (
    <>
      <Card
        className="mx-5 mb-5 h-243 w-1/3 rounded-[20px] bg-[#404040] fade-animation"
        shadow="sm"
        isPressable
      >
        <CardHeader className="relative">
          <div className="text-[24px] text-[#c6cad6] pl-3 mt-2">
            AgiGit {result.type}
          </div>
        </CardHeader>
        <CardBody className="px-3 py-0">
          <div className="p-5 flex flex-row text-center">
            <div className="basis-1/3">
              <div className="text-[28px] font-bold mb-3">You</div>
              <Code>{formatAddress(extraData.from)}</Code>
            </div>
            <div className="basis-1/3">
              <div className="text-[18px] text-[#4faaeb]">
                {extraData.gas} USDT
              </div>
              <div className="text-[28px] m-2 h-[27px]">
                <Image className="m-auto" src={dottedArrow} alt="dottedArrow" />
              </div>
              <div className="text-[18px] text-[#819df5]">Gas</div>
            </div>
            <div className="basis-1/3">
              <div className="text-[28px] font-bold mb-3">
                {result.relayHash || "wait"}
              </div>
              <Code>{formatAddress(extraData.to)}</Code>
            </div>
          </div>
        </CardBody>
      </Card>
    </>
  );
}
