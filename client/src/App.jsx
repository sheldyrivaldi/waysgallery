import Landing from "./pages/user/Landing";
import Home from "./pages/user/Home";
import Detail from "./pages/user/Detail";
import EditProfile from "./pages/user/EditProfile";
import Hired from "./pages/user/Hired";
import Post from "./pages/user/Post";
import DetailUser from "./pages/user/DetailUser";
import Profile from "./pages/user/Profile";
import MyOrder from "./pages/user/MyOrder";
import MyOffer from "./pages/user/MyOffer";
import Project from "./pages/user/Project";
import ProjectRevision from "./pages/user/ProjectRevision";
import ViewProject from "./pages/user/ViewProject";
import Withdrawal from "./pages/admin/Withdrawal";
import WithdrawalUser from "./pages/user/WithdrawalUser";

import { PrivateRouteUser, PrivateRouteAdmin, PrivateRouteLogin } from "./components/PrivateRoutes";
import { API, setAuthToken } from "./config/api";
import { useContext, useState, useEffect } from "react";
import { UserContext } from "./context/Context";
import { Routes, Route, useNavigate } from "react-router-dom";

function App() {
  const [state, dispatch] = useContext(UserContext);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  const checkUser = async () => {
    try {
      const response = await API.get("/check-auth");

      let payload = response.data.data.user;
      payload.token = localStorage.token;

      dispatch({
        type: "USER_SUCCESS",
        payload,
      });
      setIsLoading(false);
    } catch (err) {
      console.log("Check user failed : ", err);
      dispatch({
        type: "AUTH_ERROR",
      });
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (localStorage.token) {
      setAuthToken(localStorage.token);
      checkUser();
    } else {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!isLoading && !state.isLogin) {
      navigate("/");
    }
  }, [isLoading]);

  return (
    <>
      {isLoading ? null : (
        <Routes>
          <Route exact path="/" element={<Landing />} />
          <Route element={<PrivateRouteLogin />}>
            <Route element={<PrivateRouteUser />}>
              <Route exact path="/home" element={<Home />} />
              <Route exact path="/profile" element={<Profile />} />
              <Route exact path="/profile/:id" element={<DetailUser />} />
              <Route exact path="/profile/edit" element={<EditProfile />} />
              <Route exact path="/post" element={<Post />} />
              <Route exact path="/post/:id" element={<Detail />} />
              <Route exact path="/post/:id/hire" element={<Hired />} />
              <Route exact path="/order" element={<MyOrder />} />
              <Route exact path="/order/:id" element={<ViewProject />} />
              <Route exact path="/offer" element={<MyOffer />} />
              <Route exact path="/project/send/:id" element={<Project />} />
              <Route exact path="/project/revision/:id" element={<ProjectRevision />} />
              <Route exact path="/withdrawal" element={<WithdrawalUser />} />
            </Route>
            <Route element={<PrivateRouteAdmin />}>
              <Route exact path="/admin" element={<Withdrawal />} />
            </Route>
          </Route>
        </Routes>
      )}
    </>
  );
}

export default App;
