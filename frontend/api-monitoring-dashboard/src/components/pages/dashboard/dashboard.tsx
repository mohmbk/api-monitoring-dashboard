import { useState , useEffect } from 'react'
import './dashboard.css'


function Dashboard() {
 
  interface Api {
  id: string;
  userId: string;
  name: string;
  url: string;
  lastStatus: string;
  lastStatusCode: number;
  lastResponseTime: number;
  lastCheckedAt: string;
  createdAt: string;
}

const [apis , setapis] = useState<Api[]>([]) ;

 useEffect(() => {

        async function fetchApis() {
          const token = localStorage.getItem("token");
          const response = await fetch("http://localhost:8080/dashboard" , {
            method : "GET",
            headers : ({
              "Authorization": `Bearer ${token}`
            })
          })

          if(!response.ok){
            alert(await response.text());
            return;
          }

          const data = await response.json();
          setapis(data) ;
        }

        fetchApis();
        const interval = setInterval(() => {
          fetchApis();
        }, 60000);
        return () => clearInterval(interval);

    }, []);

  return (
    <>
      
    </>
  )
}

export default Dashboard