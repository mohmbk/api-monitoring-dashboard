import { useState , useEffect } from 'react'
import { useParams } from 'react-router-dom'
import './history.css'


function History() {
 
interface Check {
  id: string;
  apiId: string;
  status: string;
  statusCode: number;
  responseTime: number;
  checkedAt: string;
}

const [checks, setChecks] = useState<Check[]>([]);
const {id} = useParams();
  useEffect(() => {
  fetch(`http://localhost:8080/dashboard/api/${id}/history`)
    .then((res) => res.json())
    .then((data) => setChecks(data))
    .catch((err) => console.error(err));
  }, [id]);

  return (
    <>
      {checks.map((check) => (
        <div key= {check.id}>
          <p>status : {check.status}</p>
          <p>StatusCode : {check.statusCode}</p>
          <p>responseTime : {check.responseTime}</p>
          <p>checkedAt : {check.checkedAt}</p>
        </div>
      ))}
    </>
  )
}

export default History
