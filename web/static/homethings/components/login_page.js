export const loginPageComponent = {
    props: {
        show: Boolean,
    },
    data() {
        return {
            ruleForm: {
                login: "",
                password: "",
            },
            rules: {
                login: {
                    required: true,
                    message: "Имя пользователя обязательно для заполнения",
                    trigger: "blur",
                },
                password: {
                    required: true,
                    message: "Пароль обязателен для заполнения",
                    trigger: "blur",
                },
            },
        };
    },
    methods: {
        onSubmit() {
        },
        submitForm() {
            this.$refs.ruleForm.validate((valid) => {
                if (valid) {
                    // TODO
                    alert("submit!");
                } else {
                    return false;
                }
            });
        },
    },
    template: `
    <template v-if="show">
        <el-main class="login-page">
            <el-card>
                <el-form ref="ruleForm" :model="ruleForm" :rules="rules">
                    <el-form-item prop="login">
                        <el-input type="text" v-model="ruleForm.login" placeholder="Имя пользователя"></el-input>
                    </el-form-item>
                    <el-form-item prop="password">
                        <el-input type="password" v-model="ruleForm.password" placeholder="Пароль"></el-input>
                    </el-form-item>
                    <el-button type="primary" @click="submitForm">Авторизоваться</el-button>       
                </el-form>
            </el-card>
        </el-main>
    </template>
    `
}
