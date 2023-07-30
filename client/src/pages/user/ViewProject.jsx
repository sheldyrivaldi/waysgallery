import NavbarComp from "../../components/Navbar";
import PhotoPopUp from "../../components/PhotoPopUp";

import { Modal } from "flowbite-react";
import { useQuery, useMutation } from "react-query";
import { useNavigate, useParams } from "react-router-dom";
import { useState, useEffect } from "react";
import { API, setAuthToken } from "../../config/api";

const ViewProject = () => {
  const param = useParams();
  const navigate = useNavigate();
  const [isSubmit, setIsSubmit] = useState(false);
  const [status, setStatus] = useState("");
  const [openModal, setOpenModal] = useState();
  const props = { openModal, setOpenModal };
  const [src, setSrc] = useState("");

  setAuthToken(localStorage.token);

  const { data: project } = useQuery("projectCache", async () => {
    try {
      const response = await API.get(`/project/order/${param.id}`);
      return response.data.data.Project;
    } catch (err) {
      console.log("Fetching data failed!", err);
    }
  });

  const { data: order, refetch } = useQuery("orderProjectCache", async () => {
    try {
      const response = await API.get(`/order/${param.id}`);
      return response.data.data.orders;
    } catch (err) {
      console.log("Fetching data failed!", err);
    }
  });

  const handleSubmit = useMutation(async (e) => {
    try {
      e.preventDefault();
      if (status != "") {
        const config = {
          headers: {
            "Content-Type": "application/json",
          },
        };

        const response = await API.patch(
          `/order/${param.id}`,
          {
            status: status,
          },
          config
        );

        setIsSubmit(true);
        navigate("/order");
      }
    } catch (err) {
      console.log("Accept project failed!", err);
    }
  });

  const handleShowDetailPhoto = (url) => {
    props.setOpenModal("pop-up");
    setSrc(url);
  };

  useEffect(() => {
    refetch();
    setIsSubmit(false);
  }, [isSubmit]);

  return (
    <>
      <NavbarComp />
      <form className="flex w-full py-24 px-20" onSubmit={(e) => handleSubmit.mutate(e)}>
        <div className="w-[60%] flex flex-col items-center">
          <label for="image1" className="w-[80%] cursor-pointer flex justify-center items-center h-96 rounded">
            <img onClick={() => handleShowDetailPhoto(project?.photos[0].url)} src={project?.photos[0].url} alt="image-project" className="w-full h-full object-cover" />
          </label>
          <div className="mt-5 w-[80%] flex justify-evenly">
            {project && project?.photos.length > 0
              ? project?.photos
                  .filter((item, index) => index != 0)
                  .map((item) => (
                    <label for="image2" className="w-24 cursor-pointer flex justify-center items-center h-24  rounded">
                      <img onClick={() => handleShowDetailPhoto(item.url)} src={item.url} alt="image-project" className="w-full h-full object-cover" />
                    </label>
                  ))
              : null}
          </div>
        </div>
        <div className="w-[40%] flex flex-col items-center">
          <div className="w-[70%]">
            <p className="text-black text-justify text-opacity-50">{project?.description}</p>
          </div>
          {order?.status != "Success" ? (
            <div className="flex mt-10">
              <button onClick={() => setStatus("Revision")} type="submit" className="text-sm px-3 py-1 font-bold rounded text-white bg-[#FF9900] hover:bg-[#e08905]">
                Revision
              </button>
              <button onClick={() => setStatus("Success")} type="submit" className="text-sm ml-4 px-5 py-1 font-bold rounded text-white bg-[#0ACF83] hover:bg-[#09bd78]">
                Accept
              </button>
            </div>
          ) : null}
        </div>
      </form>
      {order?.status == "Success" ? (
        <Modal dismissible show={props.openModal === "pop-up"} onClose={() => props.setOpenModal(undefined)} size="6xl">
          <PhotoPopUp photo={src} />
        </Modal>
      ) : null}
    </>
  );
};

export default ViewProject;
