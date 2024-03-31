import { Card, CardBody, CardHeader } from "@nextui-org/react";

export default function DescriptionCard(props: any) {
  const { relayHash } = props;

  return (
    <Card
      shadow="sm"
      isPressable
      className="mx-5 mb-5 h-243 w-1/4 rounded-[20px] bg-[#404040] fade-animation"
    >
      <CardHeader>
        <div className="text-[24px] text-[#c6cad6] pl-3 mt-2">
          subscribe to {relayHash}
        </div>
      </CardHeader>
      <CardBody className="px-3 py-0">
        <div className="p-5 text-[18px] text-[#c6cad6]">订阅一个新的中继器</div>
      </CardBody>
    </Card>
  );
}
