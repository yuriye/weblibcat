var vue_det = new Vue({
    el: '#form',
    data:
        {
            message: 'Запрос',
            ISBN: '',
            Author: '',
            Title: ''
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
                )
                    .then
                    (
                        response => {
                            console.log(response);
                        },
                        error => {
                            console.error(error);
                        }
                    );
            }
        }
});