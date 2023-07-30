import Navbar from "../../components/Navbar";

import { useState, useEffect } from "react";
import { Pagination } from "flowbite-react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "react-query";
import { API, setAuthToken } from "../../config/api";
const Home = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [filterBy, setFilterBy] = useState("today");
  const [totalPage, setTotalPage] = useState(1);
  const [search, setSearch] = useState("");
  const navigate = useNavigate();

  setAuthToken(localStorage.token);

  const { data: posts, refetch } =
    filterBy == "today"
      ? useQuery("postHomeCache", async () => {
          const response = await API.get("/posts");
          setTotalPage(Math.ceil(response.data.data.posts.length / 12));
          return response.data.data.posts;
        })
      : useQuery("postHomeCache", async () => {
          const response = await API.get("/posts/following");
          setTotalPage(Math.ceil(response.data.data.posts.length / 12));
          return response.data.data.posts;
        });

  const lastPostIndex = currentPage * 12;
  const firstPostIndex = lastPostIndex - 12;

  const currentPost = posts?.slice(firstPostIndex, lastPostIndex);

  const handleNavigatePostDetail = (id) => {
    navigate(`/post/${id}`);
  };

  useEffect(() => {
    refetch();
  }, [filterBy]);

  return (
    <>
      <Navbar />
      <div className="pt-24">
        {/* Filter and Search */}
        <div className="flex justify-between px-20">
          <div>
            <select
              onChange={(e) => setFilterBy(e.target.value)}
              class="bg-gray-50 border border-gray-300 text-gray-900 text-xs rounded focus:ring-blue-500 focus:border-blue-500 block w-28 px-3 py-2 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            >
              <option value="today" selected>
                Today
              </option>
              <option value="following">Following</option>
            </select>
          </div>
          <div>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg class="w-3 h-3 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                  <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z" />
                </svg>
              </div>
              <input
                type="text"
                onChange={(e) => setSearch(e.target.value)}
                id="search"
                name="seacrh"
                class="block w-full pl-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                placeholder="Search "
                required
              />
              <button
                type="button"
                class="text-white absolute right-2.5 bottom-1.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded text-xs px-2 py-1 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
                hidden
              >
                Search
              </button>
            </div>
          </div>
        </div>

        <div className="mt-8  px-20">
          <h1 className="font-medium text-xl w-full">Today's post</h1>
        </div>

        {/* Galery */}
        <div className="w-full cursor-pointer columns-4 gap-4 mt-5 px-14">
          {currentPost
            ?.filter((item) => {
              return search.toLowerCase() === "" ? item : item.title.toLowerCase().includes(search.toLocaleLowerCase());
            })
            .map((item) => (
              <div
                onClick={() => {
                  handleNavigatePostDetail(item.id);
                }}
                className="hover:scale-110 hover:transition-all"
              >
                <img className="rounded mb-4" key={item.id} src={item.photos[0]?.url} alt="post" />
              </div>
            ))}
        </div>
      </div>
      <div className="w-full mt-24 mb-14 text-center">
        <Pagination
          currentPage={currentPage}
          onPageChange={(page) => {
            setCurrentPage(page);
          }}
          totalPages={totalPage}
        />
      </div>
    </>
  );
};

export default Home;
