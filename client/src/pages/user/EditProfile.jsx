import ProfileMockup from "../../assets/profile.svg";
import NavbarComp from "../../components/Navbar";
import { useState, useContext, useEffect } from "react";
import { UserContext } from "../../context/Context";
import { useMutation, useQuery } from "react-query";
import { useNavigate } from "react-router-dom";
import { API, setAuthToken } from "../../config/api";

const EditProfile = () => {
  const navigate = useNavigate();
  const [state] = useContext(UserContext);
  const [alert, setAlert] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");
  const [form, setForm] = useState({
    avatar: "",
    banner: "",
    greeting: "",
    fullname: "",
  });

  const [image, setImage] = useState({
    avatar: "",
    banner: "",
  });

  const handleChangeImage = (e) => {
    setForm({ ...form, [e.target.name]: e.target.files });
    setImage({ ...image, [e.target.name]: URL.createObjectURL(e.target.files[0]) });
  };

  const handleChangeForm = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  setAuthToken(localStorage.token);

  const { data: user } = useQuery("userProfileCache", async () => {
    try {
      const response = await API.get(`/user/${state.user.id}`);
      return response.data.data.users;
    } catch (err) {
      console.log("Failed fetching data!", err);
    }
  });

  console.log(user);

  const handleSubmit = useMutation(async (e) => {
    e.preventDefault();

    try {
      const config = {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      };

      const formData = new FormData();
      if (form.avatar != "" && user.avatar == "") formData.set("avatar", form.avatar[0], form.avatar[0].name);
      if (form.banner != "" && user.banner == "") formData.set("banner", form.banner[0], form.banner[0].name);

      formData.set("greeting", form.greeting);
      formData.set("fullname", form.fullname);
      const response = await API.patch(`/user/${state.user.id}`, formData, config);
      navigate("/profile");
    } catch (err) {
      console.log("ini form", form);
      console.log("Create post failed!", err);
    }
  });

  const fillingForms = () => {
    if (user?.banner != "") {
      setForm((prevForm) => ({ ...prevForm, banner: user?.banner }));
      setImage((prevImage) => ({ ...prevImage, banner: user?.banner }));
    }

    if (user?.avatar != "") {
      setForm((prevForm) => ({ ...prevForm, avatar: user?.avatar }));
      setImage((prevImage) => ({ ...prevImage, avatar: user?.avatar }));
    }

    if (user?.fullname != "") {
      setForm((prevForm) => ({ ...prevForm, fullname: user?.fullname }));
    }

    if (user?.greeting != "") {
      setForm((prevForm) => ({ ...prevForm, greeting: user?.greeting }));
    }
  };

  useEffect(() => {
    fillingForms();
  }, [user]);

  return (
    <>
      <NavbarComp />
      <form onSubmit={(e) => handleSubmit.mutate(e)} className="flex w-full py-24 px-20">
        <div className="w-[60%] flex justify-center">
          <label for="banner" className="w-[100%] relative cursor-pointer flex justify-center items-center h-96 border-4 border-[#B2B2B2] border-dashed rounded">
            {image?.banner ? (
              <img className="absolute h-full w-full object-cover rounded" src={image.banner} alt="banner" />
            ) : (
              <label for="banner" className="text-center text-xl cursor-pointer font-medium">
                <span className="text-[#2FC4B2]">Upload</span> Your Best Art
              </label>
            )}
            <input onChange={handleChangeImage} type="file" className="hidden" id="banner" name="banner" />
          </label>
        </div>
        <div className="w-[40%]  flex flex-col items-center">
          <label for="avatar" className="cursor-pointer relative w-32 h-32 rounded-full overflow-hidden ">
            {image?.avatar ? <img className="absolute h-full w-full object-cover" src={image.avatar} alt="avatar" /> : <img src={ProfileMockup} alt="profile-mockup" className="w-full h-full object-cover" />}
            <input onChange={handleChangeImage} type="file" className="hidden" id="avatar" name="avatar" />
          </label>
          <div className="w-[70%] my-8">
            <input onChange={handleChangeForm} type="text" name="greeting" className="text-lg py-2 px-3 w-full rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Greeting" value={form.greeting} />
            <input onChange={handleChangeForm} type="text" name="fullname" className="text-lg py-2 px-3 mt-4 w-full rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Fullname" value={form.fullname} />
          </div>
          <button type="submit" className="px-8 py-2 font-bold text-sm bg-[#2FC4B2] hover:bg-[#1a9b8c] text-white rounded">
            Save
          </button>
        </div>
      </form>
    </>
  );
};

export default EditProfile;
