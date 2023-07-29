import NavbarComp from "../../components/Navbar";
import Upload from "../../assets/upload.svg";
import Plus from "../../assets/plus.svg";

import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation } from "react-query";
import { API, setAuthToken } from "../../config/api";
import { CreatePostFailed } from "../../components/alert/CreatePostAlert";
const ProjectRevision = () => {
  const param = useParams();
  const navigate = useNavigate();
  const [alert, setAlert] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");
  const [image, setImage] = useState({
    image1: "",
    image2: "",
    image3: "",
    image4: "",
    image5: "",
  });

  const [form, setForm] = useState({
    image1: "",
    image2: "",
    image3: "",
    image4: "",
    image5: "",
    description: "",
  });

  setAuthToken(localStorage.token);

  const handleChangeImage = (e) => {
    setForm({ ...form, [e.target.name]: e.target.files });
    setImage({ ...image, [e.target.name]: URL.createObjectURL(e.target.files[0]) });
  };

  const handleChangeForm = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = useMutation(async (e) => {
    try {
      e.preventDefault();

      const config = {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      };

      const formData = new FormData();
      if (form.image1 != "") formData.set("image1", form.image1[0], form.image1[0].name);
      if (form.image2 != "") formData.set("image1", form.image2[0], form.image2[0].name);
      if (form.image3 != "") formData.set("image1", form.image3[0], form.image3[0].name);
      if (form.image4 != "") formData.set("image1", form.image4[0], form.image4[0].name);
      if (form.image5 != "") formData.set("image1", form.image5[0], form.image5[0].name);
      formData.set("description", form.description);

      if (form.image1 == "" && form.image2 == "" && form.image3 == "" && form.image4 == "" && form.image5 == "") {
        setAlert(true);
        setAlertMessage("Create project failed! Please upload at least 1 image.");
      } else if (form.description == "" && form.title == "") {
        setAlert(true);
        setAlertMessage("Create project failed! Title and description cannot be empty!");
      } else {
        const response = await API.post(`/project/${param.id}`, formData, config);
        navigate("/offer");
      }
    } catch (err) {
      setAlert(true);
      setAlertMessage("Server internal error!");
      console.log("Create project failed!", err);
    }
  });

  return (
    <>
      <NavbarComp />
      <div className="absolute z-[30] top-24 left-[33%]">{alert ? <CreatePostFailed message={alertMessage} /> : null}</div>
      <form onSubmit={(e) => handleSubmit.mutate(e)} className="flex w-full py-24 px-20">
        <div className="w-[60%] flex flex-col items-center">
          <label for="image1" className="w-[80%] relative file:cursor-pointer flex justify-center items-center h-96 border-4 border-[#B2B2B2] border-dashed rounded">
            <label for="image1" className="text-center text-xl cursor-pointer font-medium">
              <img src={Upload} alt="image1" className="w-52 h-52" />
              <span className="text-[#2FC4B2] mt-5">Browse</span> to choose a file
            </label>
            <input onChange={handleChangeImage} type="file" className="hidden" id="image1" name="image1" />
            {image?.image1 ? <img className="absolute h-full w-full object-cover" src={image?.image1} alt="selected-image" /> : null}
          </label>
          <div className="mt-5 w-[80%] flex justify-between">
            <label for="image2" className="w-24 relative cursor-pointer flex justify-center items-center h-24 border-4 border-[#B2B2B2] border-dashed rounded">
              <img src={Plus} alt="image2" className="w-14 h-14" />
              <input onChange={handleChangeImage} type="file" className="hidden" id="image2" name="image2" />
              {image?.image2 ? <img className="absolute h-full w-full object-cover" src={image?.image2} alt="selected-image" /> : null}
            </label>
            <label for="image3" className="w-24 relative cursor-pointer flex justify-center items-center h-24 border-4 border-[#B2B2B2] border-dashed rounded">
              <img src={Plus} alt="image3" className="w-14 h-14" />
              <input onChange={handleChangeImage} type="file" className="hidden" id="image3" name="image3" />
              {image?.image3 ? <img className="absolute h-full w-full object-cover" src={image?.image3} alt="selected-image" /> : null}
            </label>
            <label for="image4" className="w-24 relative cursor-pointer flex justify-center items-center h-24 border-4 border-[#B2B2B2] border-dashed rounded">
              <img src={Plus} alt="image4" className="w-14 h-14" />
              <input onChange={handleChangeImage} type="file" className="hidden" id="image4" name="image4" />
              {image?.image4 ? <img className="absolute h-full w-full object-cover" src={image?.image4} alt="selected-image" /> : null}
            </label>
            <label for="image5" className="w-24 relative cursor-pointer flex justify-center items-center h-24 border-4 border-[#B2B2B2] border-dashed rounded">
              <img src={Plus} alt="image5" className="w-14 h-14" />
              <input onChange={handleChangeImage} type="file" className="hidden" id="image5" name="image5" />
              {image?.image5 ? <img className="absolute h-full w-full object-cover" src={image?.image5} alt="selected-image" /> : null}
            </label>
          </div>
        </div>
        <div className="w-[40%] flex flex-col items-center">
          <div className="w-[70%]">
            <textarea onChange={handleChangeForm} name="description" className="text-lg py-2 px-3 w-full h-40 rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Description"></textarea>
          </div>
          <div className="mx-auto mt-10">
            <button type="submit" className="py-2 px-4 ml-5 rounded-md text-sm font-medium hover:bg-[#1a9b8c] text-white bg-[#2FC4B2]">
              Send Project
            </button>
          </div>
        </div>
      </form>
    </>
  );
};

export default ProjectRevision;
