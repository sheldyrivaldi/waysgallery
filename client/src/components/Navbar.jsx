import Logo from "../assets/logo.svg";
import User from "../assets/user.png";
import Dropdown from "./Dropdown";

import { Avatar } from "flowbite-react";
import { useNavigate } from "react-router-dom";
import { useContext, useState } from "react";
import { UserContext } from "../context/Context";
import { useQuery } from "react-query";
import { API, setAuthToken } from "../config/api";

const NavbarComp = () => {
  const [showDropdown, setShowDropdown] = useState(false);
  const [state] = useContext(UserContext);
  const navigate = useNavigate();

  setAuthToken(localStorage.token);

  const { data: user } = useQuery("navbarCache", async () => {
    const response = await API.get(`/user/${state.user.id}`);
    return response.data.data.users;
  });

  const handleShowDropdown = () => {
    setShowDropdown(!showDropdown);
  };

  const handleDropDropdown = () => {
    setShowDropdown(false);
  };

  const handleNavigateUploadPost = () => {
    navigate("/post");
  };

  const handleClickNavigateHome = () => {
    if (state.user.role == "admin") {
      navigate("/admin");
    } else {
      navigate("/home");
    }
  };
  return (
    <>
      <div className="flex justify-between items-center py-3 px-20 border-b fixed w-full z-[30] top-0 left-0 bg-white border-[#E1E1E1]">
        <div onClick={handleClickNavigateHome} className="w-20 cursor-pointer">
          <img src={Logo} alt="logo" />
        </div>
        <div className="flex items-center">
          <button onClick={handleNavigateUploadPost} type="button" className="px-6 py-1 text-sm font-bold text-white mr-8 rounded hover:bg-[#1a9b8c] bg-[#2FC4B2]">
            Upload
          </button>
          <div onClick={handleShowDropdown} className="w-12 cursor-pointer rounded-full">
            {user?.avatar ? <Avatar img={user.avatar} alt="avatar" rounded /> : <Avatar className="m-0 p-0" img={User} alt="logo" rounded />}
          </div>
        </div>
      </div>
      {showDropdown ? (
        <div className="modal-dropdown activate">
          <Dropdown dropDropdown={handleDropDropdown} />
        </div>
      ) : (
        <div className="modal-dropdown">
          <Dropdown dropDropdown={handleDropDropdown} />
        </div>
      )}
    </>
  );
};

export default NavbarComp;
