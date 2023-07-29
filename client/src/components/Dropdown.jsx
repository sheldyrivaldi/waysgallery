import User from "../assets/user.svg";
import Transaction from "../assets/transaction.svg";
import Logout from "../assets/logout.svg";
import Withdrawal from "../assets/withdraw.png";

import { useNavigate } from "react-router-dom";
import { useContext } from "react";
import { UserContext } from "../context/Context";
const Dropdown = ({ dropDropdown }) => {
  const handleDropDropdown = () => {
    dropDropdown();
  };
  const navigate = useNavigate();
  const [_, dispatch] = useContext(UserContext);

  const handleNavigateProfile = () => {
    navigate("/profile");
  };

  const handleNavigateOrder = () => {
    navigate("/order");
  };

  const handleNavigateWithdraw = () => {
    navigate("/withdrawal");
  };

  const handleLogout = () => {
    dispatch({
      type: "LOGOUT",
    });
  };
  return (
    <>
      <div onClick={handleDropDropdown} className="w-screen h-screen inset-0 absolute m-0 p-0 flex justify-center bg-transparent"></div>
      <div className="mt-20 w-36 h-58 pt-1 pl-2 flex flex-col absolute right-20 z-30 rounded bg-white shadow-[rgba(50,50,93,0.25)_0px_6px_12px_-2px,_rgba(0,0,0,0.3)_0px_3px_7px_-3px]">
        <div class="segitiga-dropdown absolute -top-3 right-2"></div>
        <div onClick={handleNavigateProfile} className="flex items-center py-3 cursor-pointer">
          <img className="w-8" src={User} alt="user-icon" />
          <h3 className="ml-3">Profile</h3>
        </div>
        <div onClick={handleNavigateOrder} className="flex items-center py-3 cursor-pointer">
          <img className="w-8" src={Transaction} alt="transaction-icon" />
          <h3 className="ml-3">Order</h3>
        </div>
        <div onClick={handleNavigateWithdraw} className="flex items-center py-3 cursor-pointer">
          <img className="w-8" src={Withdrawal} alt="withdrawal-icon" />
          <h3 className="ml-3">Withdrawal</h3>
        </div>
        <div className="border-t-[3px] -ml-2 border-[#A8A8A8]"></div>
        <div onClick={handleLogout} className="flex items-center ml-1 py-3 cursor-pointer">
          <img className="w-8" src={Logout} alt="logout-icon" />
          <h3 className="ml-3">Logout</h3>
        </div>
      </div>
    </>
  );
};

export default Dropdown;
