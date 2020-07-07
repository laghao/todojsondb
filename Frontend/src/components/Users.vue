    <template>
      <div class="container">
        <h3>Users:</h3>
        <table class="table">
          <thead>
            <tr>
              <th scope="col">Id</th>
              <th scope="col">Name</th>
              <th scope="col">Email</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" v-bind:key="user.id"> 
              <th scope="row">{{user.title}}</th>
              <td>{{user.completed}}</td>
              <td>{{user.created_at}}</td>
            </tr>
          </tbody>
        </table> 
      </div> 
    </template>

<script>
  import axios from 'axios';

  export default {
    name: 'Users',
    data() {
      return {
        users: null,
      };
    },
    created: function() {
      axios
        .get('https://raw.githubusercontent.com/laghao/todojsondb/ca0aa95fbd384f2a0b1196ed82c9aa00c7824515/db.json')
        .then(res => {
          this.users = res.data.todos;
        })
    }
  }
</script>

<style>
  h3 {
    margin-bottom: 5%;
  }
</style>