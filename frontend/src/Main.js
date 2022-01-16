import { useHistory, Link } from "react-router-dom";

function Main({ data }) {
  const history = useHistory();

  return (
    <div>
      <h1>Inventory</h1>
      <button onClick={() => history.push("/create_group")}>Create Group</button>
      <button onClick={() => history.push("/create_item")}>Create Item</button>
      {data.map(g => (
        <ul key={g.id}>
          <Link to={`/groups/${g.id}`}><strong style={{ fontSize: 20 }}>{g.name}</strong></Link>
          {g.items?.map(i => <Link to={`/items/${i.id}`} key={i.id}><li>{i.name}</li></Link>)}
        </ul>
      ))}
    </div>
  );
}

export default Main;