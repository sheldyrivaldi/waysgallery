import NavbarComp from "../../components/Navbar";
import Logo from "../../assets/logo.svg";
import { Avatar } from "flowbite-react";
import Rectangle from "../../assets/rectangle.svg";
import Rectangle3 from "../../assets/rectangle-3.png";
import Rectangle6 from "../../assets/rectangle-6.png";
const DetailUser = () => {
  return (
    <>
      <NavbarComp />
      <div className="px-20 py-24 overflow-hidden">
        <div className="flex">
          <div className="w-[40%] h-[70vh] flex flex-col items-start justify-center">
            <Avatar size="lg" alt="profile" img={Logo} rounded />
            <h2 className="text-lg mt-4 font-bold">Gerald</h2>
            <h2 className="text-5xl font-bold mt-4">Hey, Thanks for Looking</h2>
            <div className="mt-8">
              <button type="button" className="mr-5 font-medium text-sm py-1 px-5 rounded text-black bg-[#E7E7E7] hover:bg-[#bdbbbb]">
                Follow
              </button>
              <button type="button" className="font-medium text-sm py-1 px-7 rounded text-white bg-[#2FC4B2] hover:bg-[#1a9b8c]">
                Hire
              </button>
            </div>
          </div>
          <div className="w-[60%] relative">
            <img className="absolute w-1/3 -top-5 -right-20" src={Rectangle} alt="rectangle" />
            <div className="w-full h-full flex justify-center items-center">
              <img className="w-[70%] -mr-20 z-20" src={Rectangle3} alt="banner" />
            </div>
          </div>
        </div>
        <div className="mt-5">
          <h3 className="font-medium text-lg">Gerald Works</h3>
          <div className="w-full mt-5 grid grid-cols-4 gap-4">
            <img src={Rectangle6} alt="project" />
          </div>
        </div>
      </div>
    </>
  );
};

export default DetailUser;
