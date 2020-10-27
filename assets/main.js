var vue_det = new Vue({
    el: '#form',
    data:
        {
            message: 'Запрос',
            ISDN: '',
            Author: '',
            Title: ''
        },
    methods:
        {
            find: function () {
                this.$http.post
                (
                    '/find',
                    {
                        isdn: this.ISDN,
                        author: this.Author,
                        title: this.Title
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