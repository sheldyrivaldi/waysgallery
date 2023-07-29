import Hero from "../../assets/hero.svg";
import Corner1 from "../../assets/corner1.svg";
import Corner2 from "../../assets/corner2.svg";
import Corner3 from "../../assets/corner3.svg";
import Logo from "../../assets/logo.svg";
import Login from "../../components/auth/Login";
import Register from "../../components/auth/Register";

import { useState, useEffect } from "react";
import { Modal } from "flowbite-react";
import { useNavigate } from "react-router-dom";
import { useContext } from "react";
import { UserContext } from "../../context/Context";

const Landing = () => {
  const [openModalLogin, setOpenModalLogin] = useState();
  const propsLogin = { openModalLogin, setOpenModalLogin };
  const [openModalRegister, setOpenModalRegister] = useState();
  const propsRegister = { openModalRegister, setOpenModalRegister };

  const navigate = useNavigate();
  const [state] = useContext(UserContext);

  const checkLanding = () => {
    if (state.user.role == "admin") {
      navigate("/admin");
    } else {
      navigate("/home");
    }
  };

  useEffect(() => {
    if (state.isLogin == true) {
      checkLanding();
    }
  }, []);

  return (
    <>
      <div className="w-sceen h-screen flex overflow-hidden">
        <div className="w-[40%] h-full ">
          <img className="w-[350px] absolute top-0 left-0" src={Corner2} alt="corner2" />
          <img className="w-64 absolute bottom-0 left-0" src={Corner3} alt="corner3" />
          <div className="w-[90%] relative top-40 left-24">
            <img className="w-[22rem]" src={Logo} alt="logo" />
            <h2 className="font-medium text-2xl">show your work to inpire everyone</h2>
            <h3 className="font-light text-sm">Ways Exhibition is a website design creators gather to share their work with other creators</h3>
            <div className="mt-6">
              <button onClick={() => propsRegister.setOpenModalRegister("dismissible")} type="button" className="py-2 px-6 rounded font-bold hover:bg-[#1a9b8c] text-white bg-[#2FC4B2]">
                Join Now
              </button>
              <button onClick={() => propsLogin.setOpenModalLogin("dismissible")} type="button" className="py-2 px-11 ml-5 rounded font-bold text-black hover:bg-[#bdbbbb] bg-[#E7E7E7]">
                login
              </button>
            </div>
          </div>
        </div>
        <div className="w-[60%] h-full">
          <div className="w-full h-full flex justify-center">
            <img className="w-[75%] relative" src={Hero} alt="hero" />
            <img className="w-56 absolute bottom-0 right-0" src={Corner1} alt="corner1" />
          </div>
        </div>
      </div>
      <Modal dismissible show={propsLogin.openModalLogin === "dismissible"} onClose={() => propsLogin.setOpenModalLogin(undefined)} size="md">
        <Login showRegister={() => propsRegister.setOpenModalRegister("dismissible")} dropLogin={() => propsLogin.setOpenModalLogin(undefined)} />
      </Modal>
      <Modal dismissible show={propsRegister.openModalRegister === "dismissible"} onClose={() => propsRegister.setOpenModalRegister(undefined)} size="md">
        <Register showLogin={() => propsLogin.setOpenModalLogin("dismissible")} dropRegister={() => propsRegister.setOpenModalRegister(undefined)} />
      </Modal>
    </>
  );
};

export default Landing;
