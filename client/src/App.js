import React, { useMemo, useState, useEffect } from "react";
import axios from "axios";

import Table from "./Table";
import "./App.css";

function App() {
  const columns = useMemo(
    () => [
      {
        Header: "Sneakers",
        columns: [
          {
            Header: "Name",
            accessor: "NAME"
          },
          {
            Header: "Price",
            accessor: "PRICE"
          },
          {
            Header: "Link",
            accessor: "LINK",
            Cell: ({ cell: { value } }) => (
              <a href={value} target="_blank" rel="noopener noreferrer">
                Link
              </a>
            )
          }
        ]
      }
    ],
    []
  );

  const [data, setData] = useState([]);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [host, setHost] = useState("");
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    if (loggedIn) {
      fetchSneakerData();
    }
  }, [loggedIn]);

  const fetchSneakerData = async () => {
    try {
      const result = await axios.get(
        "http://localhost:8483/protected?action=sneaker_db_data"
      );
      setData(result.data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  const handleLogin = async () => {
    try {
      const response = await axios.post("http://localhost:8483/connect", {
        username,
        password,
        host,
      });

      if (response.status === 200) {
        setLoggedIn(true);
      } else {
        console.error("Login failed");
      }
    } catch (error) {
      console.error("Error during login:", error);
    }
  };

  return (
    <div className="App">
      {loggedIn ? (
        <Table columns={columns} data={data} />
      ) : (
        <div className="login-container">
          <h2>Login</h2>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <input
            type="text"
            placeholder="Host"
            value={host}
            onChange={(e) => setHost(e.target.value)}
          />
          <button onClick={handleLogin}>Login</button>
        </div>
      )}
    </div>
  );
}

export default App;
