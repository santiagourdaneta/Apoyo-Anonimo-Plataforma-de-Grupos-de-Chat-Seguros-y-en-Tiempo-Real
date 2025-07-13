// apoyo-anonimo/frontend/src/components/ChatRoom.jsx
import React, { useState, useEffect, useRef } from 'react';
import { getGroupById, getMessagesByGroup, sendMessage } from '../services/api';
import { useParams, useNavigate } from 'react-router-dom'; // Importamos useParams y useNavigate
import { Helmet } from 'react-helmet-async'; // Para el SEO dinámico

function ChatRoom() { // groupId y onLeaveGroup ya no son props
  const { id } = useParams(); // Leemos el ID del grupo de la URL
  const groupId = parseInt(id); // Convertimos el ID de string a número
  const navigate = useNavigate(); // Obtenemos la función para navegar

  const [group, setGroup] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [username, setUsername] = useState(() => {
    // Generar un nombre anónimo al cargar el componente de forma segura
    let randomValue;
    do {
      randomValue = crypto.getRandomValues(new Uint32Array(1))[0];
    } while (randomValue >= Math.floor((2 ** 32) / 10000) * 10000);
    randomValue %= 10000;
    return `Anónimo${randomValue}`;
  });
  const [isSending, setIsSending] = useState(false); // Estado para deshabilitar botón al enviar
  const messagesEndRef = useRef(null);
  const ws = useRef(null);

  useEffect(() => {
    const fetchGroupAndMessages = async () => {
      try {
        const groupData = await getGroupById(groupId);
        setGroup(groupData);
        const initialMessages = await getMessagesByGroup(groupId);
        setMessages(initialMessages);
      } catch (err) {
        console.error('Error al cargar el grupo o mensajes:', err);
        alert('No se pudo cargar el grupo o los mensajes. Volviendo a la lista.');
        navigate('/'); // Volver a la lista de grupos si falla la carga
      }
    };

    fetchGroupAndMessages();

    // Configurar WebSocket para mensajes en tiempo real
    const wsUrl = `ws://localhost:8000/api/ws/groups/${groupId}`;
    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log('Conectado al WebSocket del grupo:', groupId);
    };

    ws.current.onmessage = (event) => {
      const receivedMessage = JSON.parse(event.data);
      setMessages((prevMessages) => [...prevMessages, receivedMessage]);
    };

    ws.current.onclose = () => {
      console.log('Desconectado del WebSocket del grupo:', groupId);
    };

    ws.current.onerror = (err) => {
      console.error('Error en WebSocket:', err);
    };

    // Limpiar la conexión WebSocket al salir del componente
    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [groupId, navigate]); // Dependencias: groupId y navigate

  // Hacer scroll al final de los mensajes cuando se actualizan
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (newMessage.trim() === '') {
      console.log("DEBUG: Mensaje vacío, no se envía.");
      return;
    }
    setIsSending(true); // Deshabilitar botón

    try {
      const messageData = {
        group_id: groupId, // Asegúrate que la clave JSON sea snake_case
        username: username,
        content: newMessage,
      };

      console.log("DEBUG: Enviando mensaje:", messageData);

      await sendMessage(messageData);
      setNewMessage('');
    } catch (err) {
      console.error('Error al enviar mensaje:', err);
      alert('No se pudo enviar el mensaje. Intenta de nuevo.');
    } finally {
      setIsSending(false); // Habilitar botón de nuevo
    }
  };

  if (!group) return <div className="card">Cargando sala de chat...</div>;

  return (
    <div className="chat-room-container card">
      <Helmet> {/* Meta etiquetas dinámicas para SEO */}
        <title>{group.name} | Apoyo Anónimo</title>
        <meta name="description" content={`Únete a la conversación en el grupo de apoyo "${group.name}" sobre ${group.topic}.`} />
        <meta property="og:title" content={`${group.name} | Apoyo Anónimo`} />
        <meta property="og:description" content={`Únete a la conversación en el grupo de apoyo "${group.name}" sobre ${group.topic}.`} />
        <meta property="og:image" content="https://via.placeholder.com/1200x630/007bff/ffffff?text=GrupoDeApoyo" /> {/* Placeholder */}
        <meta property="twitter:card" content="summary_large_image" />
        <meta property="twitter:title" content={`${group.name} | Apoyo Anónimo`} />
        <meta property="twitter:description" content={`Únete a la conversación en el grupo de apoyo "${group.name}" sobre ${group.topic}.`} />
        <meta property="twitter:image" content="https://via.placeholder.com/1200x675/007bff/ffffff?text=GrupoDeApoyo" /> {/* Placeholder */}
      </Helmet>

      <button onClick={() => navigate('/')} className="leave-button"> {/* Usamos navigate */}
        ← Volver a Grupos
      </button>
      <h2>{group.name}</h2>
      <p className="group-topic">Tema: {group.topic}</p>
      {group.description && <p className="group-description">{group.description}</p>}

      <div className="messages-list">
        {messages.length === 0 ? (
          <p className="no-messages">Sé el primero en enviar un mensaje.</p>
        ) : (
          messages.map((msg) => (
            <div key={msg.id} className={`message-item ${msg.is_moderated ? 'moderated' : ''}`}>
              <span className="message-username">{msg.username}</span>
              <span className="message-timestamp">
                {new Date(msg.timestamp).toLocaleTimeString()}
              </span>
              <p className="message-content">{msg.content}</p>
              {msg.is_moderated && (
                <span className="moderated-tag">¡Mensaje moderado!</span>
              )}
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      <form onSubmit={handleSendMessage} className="message-input-form">
        <input
          type="text"
          placeholder={`Escribe tu mensaje como ${username}...`}
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          maxLength="500"
          aria-label="Escribe tu mensaje"
        />
        <button type="submit" disabled={isSending}>
          {isSending ? 'Enviando...' : 'Enviar'}
        </button>
      </form>
    </div>
  );
}

export default ChatRoom;
