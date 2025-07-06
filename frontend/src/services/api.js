// apoyo-anonimo/frontend/src/services/api.js

// Leemos la URL de los cerebros desde un secreto de Vite
// Si no está configurado, usamos la dirección por defecto.
const API_BASE_URL = import.meta.env.VITE_APP_API_URL || 'http://localhost:8000/api';

// Función para pedir todos los grupos
export const getGroups = async (topic = '', search = '') => {
    let url = `${API_BASE_URL}/groups`;
    const params = new URLSearchParams();
    if (topic) {
        params.append('topic', topic);
    }
    if (search) {
        params.append('search', search);
    }
    if (params.toString()) {
        url += `?${params.toString()}`;
    }

    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Error HTTP al obtener grupos: ${response.status}`);
    }
    return response.json();
};

// Función para pedir un grupo específico
export const getGroupById = async (groupId) => {
    const response = await fetch(`${API_BASE_URL}/groups/${groupId}`);
    if (!response.ok) {
        throw new Error(`Error HTTP al obtener grupo: ${response.status}`);
    }
    return response.json();
};

// Función para crear un nuevo grupo
export const createGroup = async (groupData) => {
    const response = await fetch(`${API_BASE_URL}/groups`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(groupData),
    });
    if (!response.ok) {
        throw new Error(`Error HTTP al crear grupo: ${response.status}`);
    }
    return response.json();
};

// Función para pedir los mensajes de un grupo
export const getMessagesByGroup = async (groupId) => {
    const response = await fetch(`${API_BASE_URL}/groups/${groupId}/messages`);
    if (!response.ok) {
        throw new Error(`Error HTTP al obtener mensajes: ${response.status}`);
    }
    return response.json();
};

// Función para enviar un mensaje a un grupo
export const sendMessage = async (messageData) => {
    const response = await fetch(`${API_BASE_URL}/messages`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(messageData),
    });
    if (!response.ok) {
        throw new Error(`Error HTTP al enviar mensaje: ${response.status}`);
    }
    return response.json();
};
