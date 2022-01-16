import { useState, useEffect } from "react";
import { Switch, Route, Redirect } from "react-router-dom";
import Group from "./Group";
import Item from "./Item";
import CreateGroup from "./CreateGroup";
import CreateItem from "./CreateItem";
import Main from "./Main";

function App() {
  const [notification, setNotification] = useState(null);
  const [data, setData] = useState(null);
  const [id, setId] = useState();

  useEffect(() => {
    if (!data)
      fetch('/api/groups')
        .then(response => response.json())
        .then(data => setData(data));
  }, [data, setData]);

  const notify = (notification) => {
    if (notification)
      clearTimeout(id);

    setNotification(notification);
    setId(setTimeout(() => {
      setNotification(null);
    }, 5000));
  };

  if (!data) return "Loading...";

  return (
    <div>
      {notification && <h4 style={{ color: notification.good ? "green" : "red" }}>{notification.message}</h4>}
      <Switch>
        <Route path="/groups/:id">
          <Group setData={setData} notify={notify} />
        </Route>
        <Route path="/items/:id">
          <Item data={data} setData={setData} notify={notify} />
        </Route>
        <Route path="/create_group">
          <CreateGroup setData={setData} notify={notify} />
        </Route>
        <Route path="/create_item">
          <CreateItem data={data} setData={setData} notify={notify} />
        </Route>
        <Route path="/" exact>
          <Main data={data} />
        </Route>
        <Route>
          <Redirect to="/" />
        </Route>
      </Switch>
    </div>
  );
}

export default App;
