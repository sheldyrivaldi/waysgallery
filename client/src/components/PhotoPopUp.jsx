import { saveAs } from "file-saver";

const PhotoPopUp = ({ photo }) => {
  const handleDownload = () => {
    saveAs(photo, "project");
  };
  return (
    <>
      <div className="p-10 flex flex-col justify-center items-center">
        <img className="w-full" src={photo} alt="image" />
        <button onClick={handleDownload} type="button" className="px-4 py-2 mt-5 font-bold text-sm rounded text-white bg-[#0ACF83]">
          Download
        </button>
      </div>
    </>
  );
};

export default PhotoPopUp;
