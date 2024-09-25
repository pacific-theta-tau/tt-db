import BrothersTable from './components/brothers-table'
import EventsTable from './components/events-table'
import './App.css'

function App() {
  return (
    <>
        <div className="container mx-auto py-10">
            <BrothersTable />
        </div>

        <div className="container mx-auto py-10">
            <EventsTable />
        </div>
     </>
  )
}

export default App
