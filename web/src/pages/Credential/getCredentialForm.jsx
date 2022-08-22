import { useState } from 'react'
import {
  Divider,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
} from '@mui/material'
import { toast } from 'react-toastify'
import axios from 'axios'
import QRCode from 'react-qr-code'
import validate from 'za-id-validator'
// components
import ButtonComponent from '../../components/Button'
import FormikComponent from '../../components/Formik'
import FormComponent from '../../components/Form'
import TextFieldComponent from '../../components/TextField'
import TypographyComponent from '../../components/Typography'
// utils
import { ADDRESS_ISSUER_URL } from '../../utils'

const apiURL = ADDRESS_ISSUER_URL + '/credential'

const idValidation = data => {
  const errors = {}

  if (validate(data.id_number)) {
  } else {
    errors.id_number = 'Invalid ZA ID number!'
  }

  return errors
}

const GetCredentialForm = () => {
  const [submitting, setSubmitting] = useState(false)
  const [success, setSuccess] = useState(false)
  const [data, setData] = useState([])

  const sendOffer = async data => {
    setSubmitting(true)

    await toast.promise(
      axios
        .post(apiURL, data)
        .then(response => {
          setData(response.data)
          setSuccess(true)
          toast.success('Credential request generated!')
        })
        .catch(error => {
          toast.error(error.response.data.msg)
        }),
      {
        pending: 'Generating request...',
      }
    )

    setSubmitting(false)
  }

  return (
    <>
      <TypographyComponent variant="h5">
        Get my Physical Address Credential
      </TypographyComponent>
      <Divider />
      <div style={{ marginTop: '1rem' }}>
        <FormikComponent
          initialValues={{
            id_number: '',
            forenames: '',
            surname: '',
            statement_issuer: '',
            statement_date: '',
            street_1: '',
            street_2: '',
            street_3: '',
            city: '',
            postal_code: '',
            expiry_date: '',
            email: '',
          }}
          validate={idValidation}
          onSubmit={(values, { resetForm }) => {
            sendOffer(values)
            // resetForm()
          }}
        >
          {({ values, handleChange, touched, errors }) => (
            <FormComponent>
              <div>
                <TextFieldComponent
                  error={touched.id_number && Boolean(errors.id_number)}
                  helperText={touched.id_number && errors.id_number}
                  id="id_number"
                  name="id_number"
                  value={values.id_number}
                  onChange={handleChange}
                  label="ID Number"
                  sx={{ m: '1rem' }}
                  required
                />
              </div>
              <div>
                <TextFieldComponent
                  id="forenames"
                  name="forenames"
                  value={values.forenames}
                  onChange={handleChange}
                  label="Forenames"
                  sx={{ m: '1rem' }}
                  required
                />

                <TextFieldComponent
                  id="surname"
                  name="surname"
                  value={values.surname}
                  onChange={handleChange}
                  label="Surname"
                  sx={{ m: '1rem' }}
                  required
                />
              </div>

              <div>
                <FormControl sx={{ width: '16.5rem' }}>
                  <InputLabel
                    id="statement_issuer"
                    sx={{ margin: '1rem 0 0 1rem' }}
                  >
                    Statement Issuer
                  </InputLabel>
                  <Select
                    labelId="statement_issuer"
                    id="statement_issuer"
                    name="statement_issuer"
                    value={values.statement_issuer}
                    label="Statement Issuer"
                    onChange={handleChange}
                    sx={{ m: '1rem' }}
                    required
                  >
                    <MenuItem value="Municipality">Municipality</MenuItem>
                    <MenuItem value="Department of Home Affairs">
                      Department of Home Affairs
                    </MenuItem>
                  </Select>
                </FormControl>

                <TextFieldComponent
                  id="statement_date"
                  name="statement_date"
                  type="date"
                  InputLabelProps={{ shrink: true }}
                  value={values.statement_date}
                  onChange={handleChange}
                  label="Statement Date"
                  sx={{ m: '1rem', width: '14.5rem' }}
                  required
                />
              </div>

              <div>
                <TextFieldComponent
                  id="street_1"
                  name="street_1"
                  value={values.street_1}
                  onChange={handleChange}
                  label="Street 1"
                  sx={{ m: '1rem' }}
                  required
                />

                <TextFieldComponent
                  id="street_2"
                  name="street_2"
                  value={values.street_2}
                  onChange={handleChange}
                  label="Street 2"
                  sx={{ m: '1rem' }}
                />
              </div>

              <div>
                <TextFieldComponent
                  id="street_3"
                  name="street_3"
                  value={values.street_3}
                  onChange={handleChange}
                  label="Street 3"
                  sx={{ m: '1rem' }}
                />

                <TextFieldComponent
                  id="city"
                  name="city"
                  value={values.city}
                  onChange={handleChange}
                  label="City"
                  sx={{ m: '1rem' }}
                  required
                />
              </div>

              <div>
                <TextFieldComponent
                  id="postal_code"
                  name="postal_code"
                  value={values.postal_code}
                  onChange={handleChange}
                  label="Postal Code"
                  sx={{ m: '1rem' }}
                  required
                />

                <TextFieldComponent
                  id="expiry_date"
                  name="expiry_date"
                  type="date"
                  InputLabelProps={{ shrink: true }}
                  value={values.expiry_date}
                  onChange={handleChange}
                  label="Expiry Date"
                  sx={{ m: '1rem', width: '14.5rem' }}
                  required
                />
              </div>

              <div>
                <ButtonComponent
                  variant="contained"
                  type="submit"
                  sx={{ color: '#fff', m: '1rem' }}
                  disabled={submitting}
                >
                  Submit
                </ButtonComponent>
              </div>
            </FormComponent>
          )}
        </FormikComponent>
      </div>

      {success && data.credential ? <QRCode value={data.credential} /> : ''}
    </>
  )
}

export default GetCredentialForm