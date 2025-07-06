// apoyo-anonimo/frontend/src/App.jsx
import React from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from 'react-router-dom';
import GroupList from './components/GroupList';
import ChatRoom from './components/ChatRoom';
import { HelmetProvider } from 'react-helmet-async'; // Asegúrate de que esté importado

// Componente para la página principal que lista los grupos
function Home() {
  // Home ya no necesita onSelectGroup como prop de App, GroupList lo maneja internamente
  return <GroupList />;
}

// Componente para la página de chat
function GroupChat() {
  // GroupChat ya no necesita onLeaveGroup o groupId como props de App, ChatRoom los maneja
  return <ChatRoom />;
}

function App() {
  return (
    <HelmetProvider> {/* Envuelve toda la app para Helmet */}
      <Router> {/* Envuelve las rutas */}
        <div className="app-container">
          <h1>Apoyo Anónimo</h1>
          <Routes> {/* Aquí definimos los caminos */}
            <Route path="/" element={<Home />} /> {/* Camino para la lista de grupos */}
            <Route path="/groups/:id" element={<GroupChat />} /> {/* Camino para un grupo específico */}
            {/* Puedes añadir una ruta para un 404 Not Found */}
            <Route path="*" element={
              <div className="card" style={{textAlign: 'center', padding: '50px'}}>
                <h2>Página no encontrada</h2>
                <p>Lo sentimos, la página que buscas no existe.</p>
                <button onClick={() => window.location.href = '/'}>Volver a la página principal</button>
              </div>
            } />
          </Routes>
        </div>
      </Router>
    </HelmetProvider>
  );
}

export default App;