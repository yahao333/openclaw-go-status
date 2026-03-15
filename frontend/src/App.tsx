import { Routes, Route } from 'react-router-dom'
import { Dashboard } from './pages/Dashboard'
import { Sessions } from './pages/Sessions'
import { Tasks } from './pages/Tasks'
import { Projects } from './pages/Projects'
import { Usage } from './pages/Usage'
import { Settings } from './pages/Settings'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Dashboard />} />
      <Route path="/sessions" element={<Sessions />} />
      <Route path="/tasks" element={<Tasks />} />
      <Route path="/projects" element={<Projects />} />
      <Route path="/usage" element={<Usage />} />
      <Route path="/settings" element={<Settings />} />
    </Routes>
  )
}

export default App
