import { useState } from "react";
import { useHistory, Link } from "react-router-dom";

function CreateGroup({ setData, notify }) {
  const [name, setName] = useState("");
  const history = useHistory();

  const createGroup = async event => {
    event.preventDefault();
    const options = { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ name }) };
    const response = await fetch("/api/groups", options);
    if (response.ok) {
      setData(null);
      notify({ good: true, message: `group '${name}' successfully created` });
      history.push("/");
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
        <button type="submit">Submit</button>
      </form>
    </div>
  );
}

export default CreateGroup;