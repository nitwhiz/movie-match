import { defineStore } from 'pinia';
import { User } from '../model/User';
import { computed, ref } from 'vue';
import axios from 'axios';
import { API_SERVER_BASE_URL } from '../common/Environment';

export const useUserStore = defineStore('user', () => {
  const userList = ref([] as User[]);
  const currentUser = ref(null as User | null);

  const users = computed<User[]>(() => [...userList.value]);

  const loadUsers = () => {
    axios
      .get<{ Results: User[] }>(`${API_SERVER_BASE_URL}/user`)
      .then(({ data: { Results: resultUsers } }) => {
        userList.value = resultUsers;
      });
  };

  return { currentUser, users, loadUsers };
});
