var vm = new Vue({
    el: '#form',
    data:
        {
            message: 'Запрос',
            isbn: '',
            author: '',
            title: '',
            queryResult: []
        },
    methods:
        {
            find: function () {
                this.queryResult = []
                console.log(this.isbn + "," + this.author + "," + this.title);
                this.$http.post
                (
                    '/find',
                    {
                        ISBN: this.isbn,
                        Author: this.author,
                        Title: this.title
                    }
                ).then
                    (
                        response => {
                            console.log(response.data);
                            Vue.set(vm, 'queryResult', response.data)
                        },
                        error => {
                            console.error(error);
                        }
                    );
            }
        }
});