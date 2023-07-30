import NavbarComp from "../../components/Navbar";
import { Avatar } from "flowbite-react";
import User from "../../assets/user.png";
import Rectangle from "../../assets/rectangle.svg";

import { API, setAuthToken } from "../../config/api";
import { useContext, useState } from "react";
import { useQuery } from "react-query";
import { useNavigate } from "react-router-dom";
import { UserContext } from "../../context/Context";
const Profile = () => {
  const [state] = useContext(UserContext);
  const navigate = useNavigate();

  setAuthToken(localStorage.token);
  const { data: user } = useQuery("myProfileCache", async () => {
    try {
      const response = await API.get(`/user/${state.user.id}`);
      return response.data.data.users;
    } catch (err) {
      console.log("Failed fetching data!", err);
    }
  });

  const handleNavigateEditProfile = () => {
    navigate("/profile/edit");
  };

  const handleNavigatePostDetail = (id) => {
    navigate(`/post/${id}`);
  };
  return (
    <>
      <NavbarComp />
      <div className="px-20 py-24 overflow-hidden">
        <div className="flex">
          <div className="w-[40%] h-[70vh] flex flex-col items-start justify-center">
            {user?.avatar ? (
              <div className="ml-5">
                <Avatar size="lg" alt="profile" img={user.avatar} rounded />{" "}
              </div>
            ) : (
              <div className="ml-6">
                <Avatar size="lg" alt="profile" img={User} rounded />
              </div>
            )}
            <h2 className="text-lg mt-4 ml-1.5 font-bold">{user?.fullname}</h2>
            <h2 className="text-sm  ml-1.5 font-medium">{user?.followers.length} Followers</h2>
            <h2 className="text-5xl font-bold mt-4">{user?.greeting}</h2>
            <div className="mt-8">
              <button onClick={handleNavigateEditProfile} type="button" className="font-medium text-sm py-2 px-7 rounded text-white bg-[#2FC4B2] hover:bg-[#1a9b8c]">
                Edit Profile
              </button>
            </div>
          </div>
          <div className="w-[60%] relative">
            <img className="absolute w-1/3 -top-5 -right-20" src={Rectangle} alt="rectangle" />
            <div className="w-full h-full flex justify-center items-center">
              {user?.banner ? <img className="w-full h-96 absolute bottom-20 rounded-lg" src={user?.banner} alt="banner" /> : <div className="w-full h-96 absolute bottom-20 bg-black bg-opacity-30 rounded-lg"></div>}
            </div>
          </div>
        </div>
        <div className="mt-5">
          <h3 className="font-bold text-lg">My Works</h3>

          {user?.post && user?.post.length > 0 ? (
            <div className="w-full mt-5 grid grid-cols-4 gap-4">
              {user?.post.map((item) => {
                return (
                  <div onClick={() => handleNavigatePostDetail(item.id)} className="hover:scale-110 hover:transition-all">
                    <img className="rounded" key={item.id} src={item.photos[0]?.url} alt="project" />
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center my-20">
              <h3 className="text-lg font-medium"> No Post</h3>
            </div>
          )}
        </div>
      </div>
    </>
  );
};

export default Profile;
