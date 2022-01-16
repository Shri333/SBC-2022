import { useState, useEffect } from "react";
import { useParams, useHistory, Redirect, Link } from "react-router-dom";
import { toInt } from "./utils";

function Item({ data, setData, notify }) {
  const [item, setItem] = useState(null);
  const [visible, setVisible] = useState(false);
  const [name, setName] = useState("");
  const [count, setCount] = useState(0);
  const [groupId, setGroupId] = useState(0);
  const { id } = useParams();
  const history = useHistory();

  useEffect(() => {
    let mounted = true;
    if (item === null && mounted)
      fetch(`/api/items/${id}`)
        .then(response => {
          if (response.ok)
            return response.json();
          setItem(undefined);
        })
        .then(data => {
          setItem(data);
          setName(data.name);
          setCount(data.count);
          setGroupId(data.group?.id);
        })
    return () => { mounted = false; }
  }, [id, item]);

  const updateItem = async event => {
    event.preventDefault();
    let body;
    if (groupId > 0)
      body = JSON.stringify({ name, count: toInt(count), groupId: toInt(groupId) });
    else
      body = JSON.stringify({ name, count: toInt(count) });
    const options = { method: "PUT", headers: {"content-type": "application/json"}, body };
    const response = await fetch(`/api/items/${id}`, options);
    if (response.ok) {
      const data = await response.json();
      setItem(data);
      setName(data.name);
      setCount(data.count);
      setGroupId(data.group?.id);
      setData(null);
      notify({ good: true, message: "item successfully updated" });
    } else {
      const message = await response.text();
      notify({ good: false, message });
    }
  };

  const deleteItem = async () => {
    await fetch(`/api/items/${id}`, { method: "DELETE" });
    setData(null);
    notify({ good: true, message: `item '${name}' successfully deleted` });
    history.push("/");
  };

  if (item === null) return "Loading...";

  if (item === undefined) return <Redirect to="/" />;

  return (
    <div>
      <Link to="/"><button>back</button></Link>
      <h3>{item.name}</h3>
      <p>count: {item.count}</p>
      <p>group: {item.group?.name}</p>
      {visible 
        ? (
          <form onSubmit={updateItem}>
            <div>
              Name:
              <input value={name} onChange={({ target }) => setName(target.value)} />
            </div>
            <div>
              Count:
              <input value={count} onChange={({ target }) => setCount(target.value)} />
            </div>
            <div>
              Group:
              <select value={groupId} onChange={({ target }) => setGroupId(target.value)}>
                {data.map(g => <option key={g.id} value={g.id}>{g.name}</option>)}
              </select>
            </div>
            <button type="submit">Submit</button>
            <button onClick={() => setVisible(false)}>Close</button>
          </form>
        )
        : <button onClick={() => setVisible(true)}>Edit Item</button>
      }
      <button onClick={deleteItem}>Delete Item</button>
    </div>
  );
}

export default Item;