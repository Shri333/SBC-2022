import { useState } from "react";
import { useHistory, Link } from "react-router-dom";
import { toInt } from "./utils";

function CreateItem({ data, setData, notify }) {
  const [name, setName] = useState("");
  const [count, setCount] = useState(0);
  const [groupId, setGroupId] = useState(0);
  const history = useHistory();

  const createGroup = async event => {
    event.preventDefault();
    let body;
    if (groupId > 0)
      body = JSON.stringify({ name, count: toInt(count), groupId: toInt(groupId) });
    else
      body = JSON.stringify({ name, count: toInt(count) });
    const options = { method: "POST", headers: { "content-type": "application/json" }, body };
    const response = await fetch("/api/items", options);
    if (response.ok) {
      const { id } = await response.json();
      setData(null);
      notify({ good: true, message: `item '${name}' successfully created` });
      history.push(`/items/${id}`);
    } else {
      const message = await response.text();
      notify({ good: false, message });
    }
  };

  return (
    <div>
      <Link to="/"><button>Back</button></Link>
      <form onSubmit={createGroup}>
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
      </form>
    </div>
  );
}

export default CreateItem;