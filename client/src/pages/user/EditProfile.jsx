import ProfileMockup from "../../assets/profile.svg";
import NavbarComp from "../../components/Navbar";
const EditProfile = () => {
  return (
    <>
      <NavbarComp />
      <form className="flex w-full py-24 px-20">
        <div className="w-[60%] flex justify-center">
          <label for="banner" className="w-[100%] cursor-pointer flex justify-center items-center h-96 border-4 border-[#B2B2B2] border-dashed rounded">
            <label for="banner" className="text-center text-xl cursor-pointer font-medium">
              <span className="text-[#2FC4B2]">Upload</span> Your Best Art
            </label>
            <input type="file" className="hidden" id="banner" name="banner" />
          </label>
        </div>
        <div className="w-[40%] flex flex-col items-center">
          <label for="avatar" className="cursor-pointer">
            <img src={ProfileMockup} alt="profile-mockup" className="w-32" />
            <input type="file" className="hidden" id="avatar" name="avatar" />
          </label>
          <div className="w-[70%] my-8">
            <input type="text" name="greeting" className="text-lg py-2 px-3 w-full rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Greeting" />
            <input type="text" name="fullname" className="text-lg py-2 px-3 mt-4 w-full rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Fullname" />
          </div>
          <button type="button" className="px-8 py-2 font-bold text-sm bg-[#2FC4B2] hover:bg-[#1a9b8c] text-white rounded">
            Save
          </button>
        </div>
      </form>
    </>
  );
};

export default EditProfile;
