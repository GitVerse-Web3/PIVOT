"use client";
import { NextUIProvider } from "@nextui-org/react";
import dynamic from "next/dynamic";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { AppContext } from "./component/wallet/appContext";

const WalletButtons = dynamic(
  () => import("./component/wallet/walletButtons"),
  {
    suspense: false,
    ssr: false,
  },
);

export default function Page() {
  const router = useRouter();
  const [loginState, setLoginState] = useState<boolean>(false);
  const ref = useRef<any>(null);

  useEffect(() => {
    if (ref.current) {
      ref.current.style.width = `${ref.current.scrollWidth + 10}px`;
    }
  }, []);

  function detectProvider(timeout = 3000) {
    return new Promise((resolve, reject) => {
      if (typeof (window as any).pontem === "undefined") {
        const timer = setTimeout(reject, timeout);
        window.addEventListener(
          "#pontemWalletInjected",
          (e) => {
            clearTimeout(timer);
            resolve((e as any).detail);
          },
          { once: true },
        );
      } else {
        resolve((window as any).pontem);
      }
    });
  }

  useEffect(() => {
    if (loginState) {
      detectProvider()
        .then((provider) => ((window as any).pontem = provider))
        .catch(() => console.log("Pontem Wallet not found"));
      void router.push("/home");
    }
  }, [loginState, router]);

  return (
    <NextUIProvider>
      <AppContext>
        <div
          className={`round-lg bg-[url("/login-bg.png")] bg-cover bg-center bg-no-repeat h-[100vh] items-center justify-center flex`}
        >
          <div className="flex flex-col items-center">
            <div className="text-[#fff] text-8xl mb-16 type-writer" ref={ref}>
              Babel - AgiGit 🤩
            </div>
            <WalletButtons setLoginState={setLoginState} />
          </div>
        </div>
      </AppContext>
    </NextUIProvider>
  );
}
