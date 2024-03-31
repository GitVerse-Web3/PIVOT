import { Navbar, NavbarBrand } from "@nextui-org/react";
import { AcmeLogo } from "../icon/AcmeLogo";

export default function Narbar() {
  return (
    <Navbar className="dark-gradient-bg rounded-[20px]">
      <NavbarBrand className="w-20">
        <AcmeLogo />
        <p className="font-bold text-[#c6cad6] text-[28px]">AGI-GIT:</p>
        <p className=" text-[#c6cad6] text-[22px] mx-5">
          Facilitating seamless collaboration in AGI training
        </p>
      </NavbarBrand>
    </Navbar>
  );
}
