export const loginPageComponent = {
    props: {
        show: Boolean,
    },
    data() {
        return {
            form: {
                dataList: [{name: null, amount: null}]
            },
            rules: {
                name: [{ required: true, message: 'name is mandatory!', trigger: 'blur' }],
                amount: [{ required: true, message: 'amounts is mandatory!', trigger: 'blur' }],
            }
        }
    },
    methods: {
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    alert("submit")
                    console.log('submit')
                } else {
                    return false;
                }
            });
        },
    },
    template: `
    <template v-if="show">
    
      <el-form :model="form" ref="form" label-width="20px">
        <el-table :data="form.dataList" border stripe>
          <el-table-column label="name">
            <template #default="scope">
              <el-form-item
                v-if="scope && scope.$index >= 0"
                label=" "
                :prop="'dataList.' + scope.$index + '.name'"
                :rules="rules.name"
              >
              <el-input v-model="scope.row.name"></el-input>
            </el-form-item>
            </template>
          </el-table-column>

          <el-table-column label="amount">
            <template #default="scope">
              <el-form-item
                v-if="scope && scope.$index >= 0"
                label=" "
                :prop="'dataList.' + scope.$index + '.amount'"
                :rules="rules.amount"
              >
              <el-input v-model="scope.row.amount"></el-input>
            </el-form-item>
            </template>
          </el-table-column>
        </el-table>
        <el-form-item>
          <el-button type="primary" @click="submitForm('form')">Submit</el-button>
        </el-form-item>
      </el-form>

    </template>
    `
}
