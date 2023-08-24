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
            accessor: "NAME" // Change to "NAME" since that's the property name in the data
          },
          {
            Header: "Price",
            accessor: "PRICE" // Change to "PRICE" since that's the property name in the data
          },
          {
            Header: "Link",
            accessor: "LINK", // Change to "LINK" since that's the property name in the data
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

  useEffect(() => {
    (async () => {
      try {
        const result = await axios.get(
          "http://localhost:8483/protected?action=sneaker_db_data"
        );
        setData(result.data);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    })();
  }, []);

  return (
    <div className="App">
      <Table columns={columns} data={data} />
    </div>
  );
}

export default App;
