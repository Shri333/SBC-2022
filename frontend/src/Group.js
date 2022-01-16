import { useState, useEffect } from "react";
import { useParams, useHistory, Redirect, Link } from "react-router-dom";

function Group({ setData, notify }) {
  const [group, setGroup] = useState(null);
  const [visible, setVisible] = useState(false);
  const [name, setName] = useState("");
  const { id } = useParams();
  const history = useHistory();

  useEffect(() => {
    let mounted = true;
    if (group === null && mounted)
      fetch(`/api/groups/${id}`)
        .then(response => {
          if (response.ok)
            return response.json();
          setGroup(undefined);
        })
        .then(data => {
          setGroup(data);
          setName(data.name);
        })
    return () => { mounted = false; }
  }, [id, group]);

  const updateGroup = async event => {
    event.preventDefault();
    const options = { method: "PUT", headers: {"content-type": "application/json"}, body: JSON.stringify({ name }) };
    const response = await fetch(`/api/groups/${id}`, options);
    if (response.ok) {
      const data = await response.json();
      setGroup(data);
      setName(data.name);
      setData(null);
      notify({ good: true, message: "group successfully updated" });
    } else {
      const message = await response.text();
      notify({ good: false, message });
    }
  };

  const deleteGroup = async () => {
    await fetch(`/api/groups/${id}`, { method: "DELETE" });
    setData(null);
    notify({ good: true, message: `group '${name}' successfully deleted` });
    history.push("/");
  };

  if (group === null) return "Loading...";

  if (group === undefined) return <Redirect to="/" />;

  return (
    <div>
      <Link to="/"><button>back</button></Link>
      <h3>{group.name}</h3>
      <ul>{group.items?.map(i => <Link to={`/items/${i.id}`} key={i.id}><li>{i.name}</li></Link>)}</ul>
      {visible 
        ? (
          <form onSubmit={updateGroup}>
            <div>
              Name:
              <input value={name} onChange={({ target }) => setName(target.value)} />
            </div>
            <button type="submit">Submit</button>
            <button onClick={() => setVisible(false)}>Close</button>
          </form>
        )
        : <button onClick={() => setVisible(true)}>Edit Group</button>
      }
      <button onClick={deleteGroup}>Delete Group</button>
    </div>
  );
}

export default Group;