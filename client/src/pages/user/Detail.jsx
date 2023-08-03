import NavbarComp from "../../components/Navbar";
import User from "../../assets/user.png";

import { Card, Avatar } from "flowbite-react";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery } from "react-query";
import { API, setAuthToken } from "../../config/api";
import { useContext, useState, useEffect } from "react";
import { UserContext } from "../../context/Context";

const Detail = () => {
  const param = useParams();
  const [state] = useContext(UserContext);
  const [isFollow, setIsFollow] = useState(false);
  const [isSubmit, setIsSubmit] = useState(false);
  const navigate = useNavigate();

  setAuthToken(localStorage.token);

  const { data: posts } = useQuery("postDetailCache", async () => {
    const response = await API.get(`/post/${param.id}`);
    return response.data.data.posts;
  });

  console.log(posts);

  const { data: user, refetch } = useQuery("userCache", async () => {
    const response = await API.get(`/user/${state.user.id}`);
    return response.data.data.users;
  });

  const followings = user?.followings;
  const postUserID = posts?.user.id;

  const handleFollow = async (e) => {
    e.preventDefault();
    try {
      const config = {
        headers: {
          "Content-Type": "application/json",
        },
      };

      const response = await API.post("/user", { following_id: String(postUserID) }, config);
      setIsSubmit(true);
    } catch (err) {
      console.log("Follow is failed!", err);
    }
  };

  const handleUnfollow = async (e) => {
    e.preventDefault();
    try {
      const config = {
        headers: {
          "Content-Type": "application/json",
        },
      };

      const response = await API.post("/user/unfollow", { following_id: String(postUserID) }, config);
      setIsSubmit(true);
    } catch (err) {
      console.log("Unfollow is failed!", err);
    }
  };

  const handleNavigateProfileUser = () => {
    try {
      navigate(`/profile/${postUserID}`);
    } catch (err) {
      console.log("Failed navgiate to profile user", err);
    }
  };

  const handleNavigateHire = () => {
    navigate(`/post/${param.id}/hire`);
  };

  useEffect(() => {
    refetch();
    setIsSubmit(false);
  }, [isSubmit]);
  return (
    <>
      <NavbarComp />
      <div className="py-24">
        <Card className="w-1/2 mx-auto">
          <div className="w-full flex justify-between items-center">
            <div className="flex items-center justify-start ">
              {posts?.user.avatar ? <Avatar alt="profile" img={posts?.user.avatar} rounded /> : <Avatar alt="profile" img={User} rounded />}
              <div onClick={handleNavigateProfileUser} className="flex cursor-pointer hover:scale-[101%] flex-col justify-center items-start ml-4">
                <h2 className="font-bold cursor-pointer text-md">{posts?.title}</h2>
                <h3 className="text-sm cursor-pointer">{posts?.user.fullname}</h3>
              </div>
            </div>
            {posts?.user.id != state.user.id ? (
              <div>
                {followings?.some((item) => item.id == postUserID) ? (
                  <button onClick={(e) => handleUnfollow(e)} type="button" className="mr-5 font-medium text-sm py-1 px-5 cursor-pointer rounded text-black bg-[#E7E7E7] hover:bg-[#bdbbbb]">
                    Unfollow
                  </button>
                ) : (
                  <button onClick={(e) => handleFollow(e)} type="button" className="mr-5 font-medium text-sm py-1 px-5 cursor-pointer rounded text-black bg-[#E7E7E7]  hover:bg-[#bdbbbb]">
                    Follow
                  </button>
                )}
                <button onClick={handleNavigateHire} type="button" className="font-medium text-sm py-1 px-7 rounded text-white bg-[#2FC4B2] hover:bg-[#1a9b8c]">
                  Hire
                </button>
              </div>
            ) : null}
          </div>
          <div class="grid gap-4 mt-4">
            <div>
              <img class="h-auto max-w-full rounded-lg" src={posts?.photos[0].url} alt="post-image" />
            </div>
            <div class="grid grid-cols-4 gap-4">
              {posts?.photos
                .filter((item, index) => index != 0)
                .map((item) => (
                  <div>
                    <img class="h-auto max-w-full rounded-lg" src={item.url} alt="post-image" />
                  </div>
                ))}
            </div>
          </div>
          <div>
            <h3 className="font-bold text-sm">
              👋 Say Hello <span className="text-[#2FC4B2]">{posts?.user.email}</span>
            </h3>
            <h3 className="text-sm mt-5">{posts?.description}</h3>
          </div>
        </Card>
      </div>
    </>
  );
};

export default Detail;
