// apoyo-anonimo/frontend/src/components/GroupList.jsx
import React, { useState, useEffect } from 'react';
import { getGroups, createGroup } from '../services/api';
import { useNavigate } from 'react-router-dom'; // Importamos useNavigate

function GroupList() { // onSelectGroup ya no es una prop
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [newGroupName, setNewGroupName] = useState('');
  const [newGroupDescription, setNewGroupDescription] = useState('');
  const [newGroupTopic, setNewGroupTopic] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [nameError, setNameError] = useState(''); // Para validación de nombre
  const [topicError, setTopicError] = useState(''); // Para validación de tema

  const navigate = useNavigate(); // Obtenemos la función para navegar

  const fetchGroups = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getGroups();
      setGroups(data);
    } catch (err) {
      setError('Error al cargar los grupos. Intenta de nuevo más tarde.');
      console.error('Error fetching groups:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateGroup = async (e) => {
    e.preventDefault();
    setNameError('');
    setTopicError('');

    let hasError = false;
    if (!newGroupName || newGroupName.length < 3) {
      setNameError('El nombre del grupo es obligatorio y debe tener al menos 3 letras.');
      hasError = true;
    }
    if (!newGroupTopic) {
      setTopicError('El tema del grupo es obligatorio.');
      hasError = true;
    }

    if (hasError) {
      return;
    }

    try {
      const createdGroup = await createGroup({
        name: newGroupName,
        description: newGroupDescription,
        topic: newGroupTopic,
      });
      setGroups([...groups, createdGroup]);
      setNewGroupName('');
      setNewGroupDescription('');
      setNewGroupTopic('');
      setShowCreateForm(false);
      alert('¡Grupo creado con éxito!'); // Feedback de éxito
    } catch (err) {
      alert('Error al crear el grupo. Revisa el nombre y el tema.');
      console.error('Error creating group:', err);
    }
  };

  // Función para ir al chat de un grupo
  const handleSelectGroup = (groupId) => {
    navigate(`/groups/${groupId}`); // Usamos navigate para cambiar la URL
  };

  useEffect(() => {
    fetchGroups();
  }, []);

  if (loading) return <div className="card">Cargando grupos...</div>;
  if (error) return <div className="card error-message">{error}</div>;

  return (
    <div className="group-list-container card">
      <h2>Grupos de Apoyo Anónimos</h2>
      <button onClick={() => setShowCreateForm(!showCreateForm)} style={{marginBottom: '15px'}}>
        {showCreateForm ? 'Cancelar' : 'Crear Nuevo Grupo'}
      </button>

      {showCreateForm && (
        <form onSubmit={handleCreateGroup} className="create-group-form">
          <h3>Crear Grupo</h3>
          <input
            type="text"
            placeholder="Nombre del Grupo (ej. 'Ansiedad y Estrés')"
            value={newGroupName}
            onChange={(e) => setNewGroupName(e.target.value)}
            required
            minLength="3"
            maxLength="50"
            className={nameError ? 'invalid' : ''}
            aria-label="Nombre del Grupo"
          />
          {nameError && <p className="error-text">{nameError}</p>}
          <input
            type="text"
            placeholder="Tema (ej. 'Ansiedad', 'Duelo', 'Adicciones')"
            value={newGroupTopic}
            onChange={(e) => setNewGroupTopic(e.target.value)}
            required
            className={topicError ? 'invalid' : ''}
            aria-label="Tema del Grupo"
          />
          {topicError && <p className="error-text">{topicError}</p>}
          <textarea
            placeholder="Descripción (opcional)"
            value={newGroupDescription}
            onChange={(e) => setNewGroupDescription(e.target.value)}
            rows="3"
            maxLength="200"
            aria-label="Descripción del Grupo"
          ></textarea>
          <button type="submit">Crear Grupo</button>
        </form>
      )}

      {groups.length === 0 ? (
        <p>No hay grupos disponibles. ¡Sé el primero en crear uno!</p>
      ) : (
        <ul className="group-items-list">
          {groups.map((group) => (
            <li key={group.id} className="group-item">
              <div className="group-info">
                <h3>{group.name}</h3>
                <p>Tema: {group.topic}</p>
                {group.description && <p>{group.description}</p>}
              </div>
              <button onClick={() => handleSelectGroup(group.id)} aria-label={`Unirse al grupo ${group.name}`}>
                Unirse
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default GroupList;
