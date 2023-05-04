import { useEffect, useState } from 'react'
import {
	Alert,
	Box,
	Button,
	FormControl,
	Grid,
	InputLabel,
	MenuItem,
	Paper,
	Select,
	TextField,
	Typography,
	useTheme,
} from '@mui/material'
import { Formik, Form } from 'formik'
import { toast } from 'react-toastify'
import axios from 'axios'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider'
import { DatePicker } from '@mui/x-date-pickers/DatePicker'
import { FileUploader } from 'react-drag-drop-files'
import { AddAPhoto } from '@mui/icons-material'

import { ADDRESS_ISSUER_URL, GMAK, LIGHT_MODE_THEME } from '../../utils/constants'

interface EnterDetailsProps {
	handleNext: () => void
}

const apiURL = ADDRESS_ISSUER_URL + '/credential'
const apiEmailURL = ADDRESS_ISSUER_URL + '/email-credential'

const fileTypes = ['JPEG', 'JPG']

const EnterDetails = ({ handleNext }: EnterDetailsProps) => {
	const theme = useTheme()

	const [submitting, setSubmitting] = useState<boolean | undefined>(false)
	const [input, setInput] = useState('')
	const [file, setFile] = useState(null)

	const sendOffer = async (values: any) => {
		setSubmitting(true)

		await toast.promise(
			axios
				.post(apiURL, values)
				.then((response: any) => {
					// @ts-ignore
					sessionStorage.setItem('credential', response.data.credential)
					toast.success('Credential generated!')
					setTimeout(() => {
						handleNext()
					}, 1000)
				})
				.catch((error: any) => {
					toast.error(error.response.data.msg)
				}),
			{
				pending: 'Generating credential...',
			}
		)
		setSubmitting(false)
	}

	const emailOffer = async (values: any) => {
		setSubmitting(true)

		await toast.promise(
			axios
				.post(apiEmailURL, values)
				.then((response: any) => {
					// @ts-ignore
					sessionStorage.setItem('credential', response.data.credential)
					toast.success('Emailed credential!')
				})
				.catch((error: any) => {
					toast.error(error.response.data.msg)
				}),
			{
				pending: 'Emailing credential...',
			}
		)
		setSubmitting(false)
	}

	const handleFileChange = (file: any) => {
		getBase64(file, (result: any) => {
			const imgBase64 = result
			setFile(imgBase64)
		})
	}

	const getBase64 = (file: any, cb: any) => {
		let reader = new FileReader()
		reader.readAsDataURL(file)
		reader.onload = function () {
			cb(reader.result)
		}
		reader.onerror = function (error) {
			console.log('Error: ', error)
		}
	}

	return (
		<Grid container>
			<Grid item xs={0} md={2} />
			<Grid item xs={12} md={8}>
				<Paper
					square
					elevation={2}
					sx={{
						p: 5,
						width: { md: '100%' },
						backgroundColor: theme.palette.mode === LIGHT_MODE_THEME ? '#fff' : '',
						borderRadius: 5,
					}}>
					<Formik
						initialValues={{
							address_line: '',
							address_line_2: '',
							city: '',
							province: '',
							postal_code: '',
							country_code: '',
							statement_date: '',
							statement_image: '',
							self_attested: '',
							email: '',
						}}
						// validate={idValidation}
						onSubmit={(values) => {
							values.email === '' ? sendOffer(values) : emailOffer(values)
						}}>
						{({ values, handleChange, setFieldValue }) => (
							<Form>
								<div>
									<TextField
										id='address_line'
										name='address_line'
										value={values.address_line}
										onChange={handleChange}
										label='Address Line'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										required
									/>
								</div>

								<div>
									<TextField
										id='address_line_2'
										name='address_line_2'
										value={values.address_line_2}
										onChange={handleChange}
										label='Address Line 2'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
									/>
								</div>

								<div>
									<TextField
										id='city'
										name='city'
										value={values.city}
										onChange={handleChange}
										label='City'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										required
									/>
								</div>

								<div>
									<TextField
										id='province'
										name='province'
										value={values.province}
										onChange={handleChange}
										label='Province'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										required
									/>
								</div>

								<div>
									<TextField
										id='postal_code'
										name='postal_code'
										value={values.postal_code}
										onChange={handleChange}
										label='Postal Code'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										required
									/>
								</div>

								<div>
									<TextField
										id='country_code'
										name='country_code'
										value={values.country_code}
										onChange={handleChange}
										label='Country Code'
										sx={{ m: '1rem', width: '70%', cursor: 'not-allowed' }}
										required
									/>
								</div>

								<div>
									<LocalizationProvider dateAdapter={AdapterDayjs}>
										<DatePicker
											label='Statement Date'
											value={values.statement_date}
											onChange={(value) => setFieldValue('statement_date', value, true)}
											renderInput={(params) => (
												<TextField
													id='statement_date'
													name='statement_date'
													{...params}
													sx={{ m: '1rem', width: '14.5rem' }}
													required
												/>
											)}
										/>
									</LocalizationProvider>

									<FormControl sx={{ width: '16.5rem' }}>
										<InputLabel id='self_attested' sx={{ margin: '1rem 0 0 1rem' }} required>
											Self Attested
										</InputLabel>
										<Select
											labelId='self_attested'
											id='self_attested'
											name='self_attested'
											value={values.self_attested}
											label='Self Attested'
											onChange={handleChange}
											required
											sx={{ m: '1rem' }}>
											<MenuItem value='N'>N</MenuItem>
											<MenuItem value='Y'>Y</MenuItem>
										</Select>
									</FormControl>
								</div>

								<FileUploader
									name='statement'
									handleChange={handleFileChange}
									types={fileTypes}
									maxSize={0.1}
									onSizeError={() => toast.error('Image is too large. Max size is 100kb.')}
									onTypeError={() => toast.error('Incorrect file type.')}>
									<Box
										sx={{
											width: '100%',
											height: '20rem',
											borderRadius: '0.5rem',
											border: '1px dashed #de5353',
											p: file ? '1rem' : '5rem',
											cursor: 'pointer',
										}}>
										{file ? (
											<>
												<img src={file} alt='statement' style={{ maxWidth: 400, maxHeight: 280 }} />
												<TextField
													id='statement_image'
													name='statement_image'
													value={(values.statement_image = file)}
													required
													sx={{ display: 'none' }}
												/>
											</>
										) : (
											<>
												<Typography variant='h5'>
													<AddAPhoto /> Upload or Drag an image of your statement
												</Typography>
												<Typography>Type: JPEG or JPG | Max Size: 100kb</Typography>
											</>
										)}
									</Box>
								</FileUploader>

								<div>
									<TextField
										id='email'
										name='email'
										type='email'
										value={values.email}
										onChange={handleChange}
										label='Email'
										sx={{ m: '1rem' }}
										// helperText='Optional: If you want to receive your credential via email.'
										helperText={
											<Alert severity='info' sx={{ backgroundColor: 'transparent' }}>
												Optional: If you want to receive your credential via email.
											</Alert>
										}
									/>
								</div>

								<div>
									<Button
										variant='contained'
										size='small'
										type='submit'
										sx={{ color: '#fff', m: '1rem' }}
										disabled={submitting}>
										Submit
									</Button>
								</div>
							</Form>
						)}
					</Formik>
				</Paper>
			</Grid>
			<Grid item xs={0} md={2} />
		</Grid>
	)
}

export default EnterDetails
