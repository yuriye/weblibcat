var vm = new Vue({
    el: '#form',
    data:
        {
            message: 'Запрос',
            ISBN: '',
            Author: '',
            Title: '',
            queryResult: []
        },
    methods:
        {
            find: function () {
                console.log(this.ISBN + "," + this.Author + "," + this.Title);
                this.$http.post
                (
                    '/find',
                    {
                        ISBN: this.ISBN,
                        Author: this.Author,
                        Title: this.Title
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