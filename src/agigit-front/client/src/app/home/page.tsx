"use client";
import { NextUIProvider, Progress } from "@nextui-org/react";
import { Contract } from "ethers";
import { useEffect, useRef, useState } from "react";
import CommandLine from "../component/commandLine";
import DescriptionCard from "../component/descriptionCard";
import Navbar from "../component/navbar";
import TransferCard from "../component/transferCard";
import { AppContext } from "../component/wallet/appContext";
import ABI from "../constant/ABI.json";
import { contractAddress } from "../constant/contractAddress";
interface commandResult {
  type: string | "add" | "remove" | "pull" | "push" | null;
  relayHash: string;
}

export default function Home() {
  function parsingCommand(command: string) {
    let result: commandResult = {
      type: null,
      relayHash: "",
    };
    const subStrList = [
      {
        substr: "agigit relay add",
        getResult: (command: string) => {
          const list = command.split(" ");
          return {
            type: "add",
            relayHash: list[list.length - 1],
          };
        },
      },
      {
        substr: "agigit relay remove",
        getResult: () => {
          return {
            type: "remove",
            relayHash: "",
          };
        },
      },
      {
        substr: "agigit pull",
        getResult: () => {
          return {
            type: "pull",
          };
        },
      },
      {
        substr: "agigit push",
        getResult: () => {
          return {
            type: "push",
          };
        },
      },
    ];
    for (let item of subStrList) {
      if (command.toLocaleLowerCase().includes(item.substr)) {
        result = {
          ...result,
          ...item.getResult(command.toLocaleLowerCase()),
        };
        console.log(result);
        break;
      } else {
        console.log("error command!");
      }
    }
    return result;
  }
  const [command, setCommand] = useState("");
  const [result, setResult] = useState<commandResult>();
  const [pullCardDom, setPullCardDom] = useState<any>();
  const [pushCardDom, setPushCardDom] = useState<any>();
  const relayHashRef = useRef("");
  const pontemProviderRef = useRef<any>(null);
  let contract = useRef<any>(null);

  useEffect(() => {
    if ((window as any).pontem) {
      pontemProviderRef.current = (window as any).pontem;
      contract.current = new Contract(
        contractAddress,
        ABI,
        pontemProviderRef.current,
      );
    }
  }, []);

  useEffect(() => {
    const value = parsingCommand(command);
    if (value.type !== "add" && value.type !== "remove") {
      value.relayHash = relayHashRef.current;
    }
    setResult(value);
  }, [command]);

  const contractOfPull = () => {
    setPullCardDom(
      <Progress
        size="sm"
        isIndeterminate
        aria-label="Loading..."
        className="max-w-md my-auto"
      />,
    );
    const timer = setTimeout(() => {
      const extraData = {
        from: "1J93t4tZ76hX9i3Qd7aR4yY3F4iGqG6z8z7gPN8oQe",
        to: "3A3gYZFopmTK1bSVLwsMQDntWQTZARfNXq",
        gas: "0.002",
      };
      setPullCardDom(
        <TransferCard key="pull" result={result} extraData={extraData} />,
      );
    }, 3000);
  };

  const contractOfPush = () => {
    setPushCardDom(
      <Progress
        size="sm"
        isIndeterminate
        aria-label="Loading..."
        className="max-w-md my-auto"
      />,
    );
    const timer = setTimeout(() => {
      const extraData = {
        from: "1J93t4tZ76hX9i3Qd7aR4yY3F4iGqG6z8z7gPN8oQe",
        to: "3A3gYZFopmTK1bSVLwsMQDntWQTZARfNXq",
        gas: "0.003821564",
      };
      setPushCardDom(
        <TransferCard key="push" result={result} extraData={extraData} />,
      );
    }, 3000);
  };
  const contractOfRaisePayment = async () => {
    try {
      const result = await pontemProviderRef.current.signAndSubmit();
    } catch (e) {
      console.log(e);
    }
  };

  useEffect(() => {
    if (result?.type === "add") {
      relayHashRef.current = result.relayHash;
    }
    if (result?.type === "remove") {
      relayHashRef.current = "";
    }
    if (result?.type === "pull") {
      contractOfPull();
    }
    if (result?.type === "push") {
      contractOfPush();
    }
  }, [result]);

  return (
    <NextUIProvider>
      <AppContext>
        <div className="flex min-h-screen flex-col p-24">
          <Navbar />
          <CommandLine setCommand={setCommand} />
          <div className="flex flex-row px-6 w-full">
            {result?.relayHash && <DescriptionCard key="add" {...result} />}
            {pullCardDom}
            {pushCardDom}
          </div>
        </div>
      </AppContext>
    </NextUIProvider>
  );
}
